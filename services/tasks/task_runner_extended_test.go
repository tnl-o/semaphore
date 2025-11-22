package tasks

import (
	"testing"
	"time"

	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/db/bolt"
	"github.com/semaphoreui/semaphore/pkg/task_logger"
	"github.com/semaphoreui/semaphore/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestPool(t *testing.T) (*bolt.BoltDb, *TaskPool) {
	t.Helper()

	store := bolt.CreateTestStore()
	pool := CreateTaskPool(
		store,
		NewMemoryTaskStateStore(),
		nil,
		&InventoryServiceMock{},
		&EncryptionServiceMock{},
		nil,
		&mockLogWriteService{},
	)

	go pool.Run()
	time.Sleep(10 * time.Millisecond)

	return store, pool
}

func waitForTaskInPool(t *testing.T, pool *TaskPool, taskID int) *TaskRunner {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if task := pool.GetTask(taskID); task != nil {
			return task
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("timed out waiting for task %d to be enqueued", taskID)
	return nil
}

type noopStore struct {
	db.Store
}

func (s *noopStore) UpdateTask(task db.Task) error {
	return nil
}

func TestTaskPool_AddTask(t *testing.T) {
	store, pool := setupTestPool(t)

	// Create project and template
	proj, err := store.CreateProject(db.Project{})
	require.NoError(t, err)

	key, err := store.CreateAccessKey(db.AccessKey{
		ProjectID: &proj.ID,
		Type:      db.AccessKeyNone,
	})
	require.NoError(t, err)

	repo, err := store.CreateRepository(db.Repository{
		ProjectID: proj.ID,
		SSHKeyID:  key.ID,
		Name:      "Test",
		GitURL:    "git@example.com:test/test",
		GitBranch: "master",
	})
	require.NoError(t, err)

	inv, err := store.CreateInventory(db.Inventory{
		ProjectID: proj.ID,
	})
	require.NoError(t, err)

	env, err := store.CreateEnvironment(db.Environment{
		ProjectID: proj.ID,
		Name:      "test",
		JSON:      `{"test": "value"}`,
	})
	require.NoError(t, err)

	tpl, err := store.CreateTemplate(db.Template{
		Name:          "Test Template",
		Playbook:      "test.yml",
		ProjectID:     proj.ID,
		RepositoryID:  repo.ID,
		InventoryID:   &inv.ID,
		EnvironmentID: &env.ID,
	})
	require.NoError(t, err)

	// Test adding a task
	task := db.Task{
		TemplateID:  tpl.ID,
		ProjectID:   proj.ID,
		Environment: `{"additional": "vars"}`,
	}

	taskRunner, err := pool.AddTask(task, nil, "", proj.ID, false)
	require.NoError(t, err)
	assert.NotNil(t, taskRunner)
	assert.Equal(t, tpl.ID, taskRunner.Template.ID)
	assert.Equal(t, proj.ID, taskRunner.Task.ProjectID)
}

func TestTaskPool_StopTask(t *testing.T) {
	store, pool := setupTestPool(t)

	// Create project and template
	proj, err := store.CreateProject(db.Project{})
	require.NoError(t, err)

	key, err := store.CreateAccessKey(db.AccessKey{
		ProjectID: &proj.ID,
		Type:      db.AccessKeyNone,
	})
	require.NoError(t, err)

	repo, err := store.CreateRepository(db.Repository{
		ProjectID: proj.ID,
		SSHKeyID:  key.ID,
		Name:      "Test",
		GitURL:    "git@example.com:test/test",
		GitBranch: "master",
	})
	require.NoError(t, err)

	inv, err := store.CreateInventory(db.Inventory{
		ProjectID: proj.ID,
	})
	require.NoError(t, err)

	env, err := store.CreateEnvironment(db.Environment{
		ProjectID: proj.ID,
		Name:      "test",
		JSON:      `{"test": "value"}`,
	})
	require.NoError(t, err)

	tpl, err := store.CreateTemplate(db.Template{
		Name:          "Test Template",
		Playbook:      "test.yml",
		ProjectID:     proj.ID,
		RepositoryID:  repo.ID,
		InventoryID:   &inv.ID,
		EnvironmentID: &env.ID,
	})
	require.NoError(t, err)

	// Add task
	task := db.Task{
		TemplateID: tpl.ID,
		ProjectID:  proj.ID,
	}

	taskRunner, err := pool.AddTask(task, nil, "", proj.ID, false)
	require.NoError(t, err)
	waitForTaskInPool(t, pool, taskRunner.Task.ID)

	// Test stopping task (graceful)
	err = pool.StopTask(taskRunner.Task, false)
	require.NoError(t, err)

	// Verify task status is set to stopping
	updatedTask := pool.GetTask(taskRunner.Task.ID)
	if updatedTask != nil {
		assert.Equal(t, task_logger.TaskStoppingStatus, updatedTask.Task.Status)
	}

	// Test force stopping task
	task2 := db.Task{
		TemplateID: tpl.ID,
		ProjectID:  proj.ID,
	}

	taskRunner2, err := pool.AddTask(task2, nil, "", proj.ID, false)
	require.NoError(t, err)
	waitForTaskInPool(t, pool, taskRunner2.Task.ID)

	err = pool.StopTask(taskRunner2.Task, true)
	require.NoError(t, err)

	// Verify task status is set to stopped
	updatedTask2 := pool.GetTask(taskRunner2.Task.ID)
	if updatedTask2 != nil {
		assert.Equal(t, task_logger.TaskStoppedStatus, updatedTask2.Task.Status)
	}
}

func TestTaskPool_StopTask_NotActive(t *testing.T) {
	store, pool := setupTestPool(t)

	// Create backing project/template so populateDetails succeeds
	proj, err := store.CreateProject(db.Project{})
	require.NoError(t, err)

	key, err := store.CreateAccessKey(db.AccessKey{
		ProjectID: &proj.ID,
		Type:      db.AccessKeyNone,
	})
	require.NoError(t, err)

	repo, err := store.CreateRepository(db.Repository{
		ProjectID: proj.ID,
		SSHKeyID:  key.ID,
		Name:      "Test",
		GitURL:    "git@example.com:test/test",
		GitBranch: "master",
	})
	require.NoError(t, err)

	inv, err := store.CreateInventory(db.Inventory{
		ProjectID: proj.ID,
	})
	require.NoError(t, err)

	env, err := store.CreateEnvironment(db.Environment{
		ProjectID: proj.ID,
		Name:      "test",
		JSON:      `{"test": "value"}`,
	})
	require.NoError(t, err)

	tpl, err := store.CreateTemplate(db.Template{
		Name:          "Test Template",
		Playbook:      "test.yml",
		ProjectID:     proj.ID,
		RepositoryID:  repo.ID,
		InventoryID:   &inv.ID,
		EnvironmentID: &env.ID,
	})
	require.NoError(t, err)

	// Create a task that doesn't exist in the pool
	task := db.Task{
		TemplateID: tpl.ID,
		ProjectID:  proj.ID,
		Status:     task_logger.TaskSuccessStatus,
	}
	createdTask, err := store.CreateTask(task, util.Config.MaxTasksPerTemplate)
	require.NoError(t, err)

	// Should not error, but mark task as stopped
	err = pool.StopTask(createdTask, false)
	require.NoError(t, err)
}

func TestTaskPool_ConfirmTask(t *testing.T) {
	store, pool := setupTestPool(t)

	// Create project and template
	proj, err := store.CreateProject(db.Project{})
	require.NoError(t, err)

	key, err := store.CreateAccessKey(db.AccessKey{
		ProjectID: &proj.ID,
		Type:      db.AccessKeyNone,
	})
	require.NoError(t, err)

	repo, err := store.CreateRepository(db.Repository{
		ProjectID: proj.ID,
		SSHKeyID:  key.ID,
		Name:      "Test",
		GitURL:    "git@example.com:test/test",
		GitBranch: "master",
	})
	require.NoError(t, err)

	inv, err := store.CreateInventory(db.Inventory{
		ProjectID: proj.ID,
	})
	require.NoError(t, err)

	env, err := store.CreateEnvironment(db.Environment{
		ProjectID: proj.ID,
		Name:      "test",
		JSON:      `{"test": "value"}`,
	})
	require.NoError(t, err)

	tpl, err := store.CreateTemplate(db.Template{
		Name:          "Test Template",
		Playbook:      "test.yml",
		ProjectID:     proj.ID,
		RepositoryID:  repo.ID,
		InventoryID:   &inv.ID,
		EnvironmentID: &env.ID,
	})
	require.NoError(t, err)

	// Add task
	task := db.Task{
		TemplateID: tpl.ID,
		ProjectID:  proj.ID,
	}

	taskRunner, err := pool.AddTask(task, nil, "", proj.ID, false)
	require.NoError(t, err)
	waitForTaskInPool(t, pool, taskRunner.Task.ID)

	// Test confirming task
	err = pool.ConfirmTask(taskRunner.Task)
	require.NoError(t, err)

	// Verify task status
	updatedTask := pool.GetTask(taskRunner.Task.ID)
	require.NotNil(t, updatedTask)
	assert.Equal(t, task_logger.TaskConfirmed, updatedTask.Task.Status)
}

func TestTaskPool_RejectTask(t *testing.T) {
	store, pool := setupTestPool(t)

	// Create project and template
	proj, err := store.CreateProject(db.Project{})
	require.NoError(t, err)

	key, err := store.CreateAccessKey(db.AccessKey{
		ProjectID: &proj.ID,
		Type:      db.AccessKeyNone,
	})
	require.NoError(t, err)

	repo, err := store.CreateRepository(db.Repository{
		ProjectID: proj.ID,
		SSHKeyID:  key.ID,
		Name:      "Test",
		GitURL:    "git@example.com:test/test",
		GitBranch: "master",
	})
	require.NoError(t, err)

	inv, err := store.CreateInventory(db.Inventory{
		ProjectID: proj.ID,
	})
	require.NoError(t, err)

	env, err := store.CreateEnvironment(db.Environment{
		ProjectID: proj.ID,
		Name:      "test",
		JSON:      `{"test": "value"}`,
	})
	require.NoError(t, err)

	tpl, err := store.CreateTemplate(db.Template{
		Name:          "Test Template",
		Playbook:      "test.yml",
		ProjectID:     proj.ID,
		RepositoryID:  repo.ID,
		InventoryID:   &inv.ID,
		EnvironmentID: &env.ID,
	})
	require.NoError(t, err)

	// Add task
	task := db.Task{
		TemplateID: tpl.ID,
		ProjectID:  proj.ID,
	}

	taskRunner, err := pool.AddTask(task, nil, "", proj.ID, false)
	require.NoError(t, err)
	waitForTaskInPool(t, pool, taskRunner.Task.ID)

	// Test rejecting task
	err = pool.RejectTask(taskRunner.Task)
	require.NoError(t, err)

	// Verify task status
	updatedTask := pool.GetTask(taskRunner.Task.ID)
	require.NotNil(t, updatedTask)
	assert.Equal(t, task_logger.TaskRejected, updatedTask.Task.Status)
}

func TestTaskPool_StopTasksByTemplate(t *testing.T) {
	store, pool := setupTestPool(t)

	// Create project and template
	proj, err := store.CreateProject(db.Project{})
	require.NoError(t, err)

	key, err := store.CreateAccessKey(db.AccessKey{
		ProjectID: &proj.ID,
		Type:      db.AccessKeyNone,
	})
	require.NoError(t, err)

	repo, err := store.CreateRepository(db.Repository{
		ProjectID: proj.ID,
		SSHKeyID:  key.ID,
		Name:      "Test",
		GitURL:    "git@example.com:test/test",
		GitBranch: "master",
	})
	require.NoError(t, err)

	inv, err := store.CreateInventory(db.Inventory{
		ProjectID: proj.ID,
	})
	require.NoError(t, err)

	env, err := store.CreateEnvironment(db.Environment{
		ProjectID: proj.ID,
		Name:      "test",
		JSON:      `{"test": "value"}`,
	})
	require.NoError(t, err)

	tpl, err := store.CreateTemplate(db.Template{
		Name:          "Test Template",
		Playbook:      "test.yml",
		ProjectID:     proj.ID,
		RepositoryID:  repo.ID,
		InventoryID:   &inv.ID,
		EnvironmentID: &env.ID,
	})
	require.NoError(t, err)

	// Add multiple tasks
	task1 := db.Task{
		TemplateID: tpl.ID,
		ProjectID:  proj.ID,
	}
	task2 := db.Task{
		TemplateID: tpl.ID,
		ProjectID:  proj.ID,
	}

	addedTask1, err := pool.AddTask(task1, nil, "", proj.ID, false)
	require.NoError(t, err)
	waitForTaskInPool(t, pool, addedTask1.Task.ID)

	addedTask2, err := pool.AddTask(task2, nil, "", proj.ID, false)
	require.NoError(t, err)
	waitForTaskInPool(t, pool, addedTask2.Task.ID)

	// Stop all tasks for this template
	pool.StopTasksByTemplate(proj.ID, tpl.ID, false)

	// Give some time for status updates
	time.Sleep(100 * time.Millisecond)

	// Verify tasks are stopped
	tasks := pool.state.QueueRange()
	for _, queuedTask := range tasks {
		if queuedTask != nil && queuedTask.Task.TemplateID == tpl.ID {
			assert.True(t, queuedTask.Task.Status == task_logger.TaskStoppingStatus || queuedTask.Task.Status == task_logger.TaskStoppedStatus)
		}
	}
}

func TestTaskRunner_SetStatus(t *testing.T) {
	pool := &TaskPool{
		store:            &noopStore{},
		state:            NewMemoryTaskStateStore(),
		logWriteService:  &mockLogWriteService{},
		inventoryService: &InventoryServiceMock{},
	}

	task := db.Task{
		ID:     1,
		Status: task_logger.TaskWaitingStatus,
	}

	taskRunner := NewTaskRunner(task, pool, "testuser", &KeyInstallerMock{})

	// Test setting status
	taskRunner.SetStatus(task_logger.TaskRunningStatus)
	assert.Equal(t, task_logger.TaskRunningStatus, taskRunner.Task.Status)

	taskRunner.SetStatus(task_logger.TaskSuccessStatus)
	assert.Equal(t, task_logger.TaskSuccessStatus, taskRunner.Task.Status)
}

func TestTaskRunner_StatusTransitions(t *testing.T) {
	pool := &TaskPool{
		store:            &noopStore{},
		state:            NewMemoryTaskStateStore(),
		logWriteService:  &mockLogWriteService{},
		inventoryService: &InventoryServiceMock{},
	}

	task := db.Task{
		ID:     1,
		Status: task_logger.TaskWaitingStatus,
	}

	taskRunner := NewTaskRunner(task, pool, "testuser", &KeyInstallerMock{})

	// Test valid status transitions
	statuses := []task_logger.TaskStatus{
		task_logger.TaskWaitingStatus,
		task_logger.TaskStartingStatus,
		task_logger.TaskRunningStatus,
		task_logger.TaskSuccessStatus,
	}

	for _, status := range statuses {
		taskRunner.SetStatus(status)
		assert.Equal(t, status, taskRunner.Task.Status)
	}
}

func TestTaskPool_GetTask(t *testing.T) {
	store, pool := setupTestPool(t)

	// Create project and template
	proj, err := store.CreateProject(db.Project{})
	require.NoError(t, err)

	key, err := store.CreateAccessKey(db.AccessKey{
		ProjectID: &proj.ID,
		Type:      db.AccessKeyNone,
	})
	require.NoError(t, err)

	repo, err := store.CreateRepository(db.Repository{
		ProjectID: proj.ID,
		SSHKeyID:  key.ID,
		Name:      "Test",
		GitURL:    "git@example.com:test/test",
		GitBranch: "master",
	})
	require.NoError(t, err)

	inv, err := store.CreateInventory(db.Inventory{
		ProjectID: proj.ID,
	})
	require.NoError(t, err)

	env, err := store.CreateEnvironment(db.Environment{
		ProjectID: proj.ID,
		Name:      "test",
		JSON:      `{"test": "value"}`,
	})
	require.NoError(t, err)

	tpl, err := store.CreateTemplate(db.Template{
		Name:          "Test Template",
		Playbook:      "test.yml",
		ProjectID:     proj.ID,
		RepositoryID:  repo.ID,
		InventoryID:   &inv.ID,
		EnvironmentID: &env.ID,
	})
	require.NoError(t, err)

	// Add task
	task := db.Task{
		TemplateID: tpl.ID,
		ProjectID:  proj.ID,
	}

	taskRunner, err := pool.AddTask(task, nil, "", proj.ID, false)
	require.NoError(t, err)

	// Test getting task
	foundTask := waitForTaskInPool(t, pool, taskRunner.Task.ID)
	assert.Equal(t, taskRunner.Task.ID, foundTask.Task.ID)
	assert.Equal(t, taskRunner.Task.TemplateID, foundTask.Task.TemplateID)

	// Test getting non-existent task
	notFoundTask := pool.GetTask(99999)
	assert.Nil(t, notFoundTask)
}

func TestTaskPool_GetRunningTasks(t *testing.T) {
	_, pool := setupTestPool(t)

	// Initially should be empty
	runningTasks := pool.GetRunningTasks()
	assert.Empty(t, runningTasks)
}

func TestTaskRunner_PopulateDetails_ErrorHandling(t *testing.T) {
	_, pool := setupTestPool(t)

	// Test with non-existent template
	task := db.Task{
		TemplateID: 99999, // Non-existent
		ProjectID:  1,
	}

	taskRunner := NewTaskRunner(task, pool, "testuser", &KeyInstallerMock{})
	err := taskRunner.populateDetails()

	require.Error(t, err)
}

func TestTaskPool_ConfirmReject_NotActive(t *testing.T) {
	_, pool := setupTestPool(t)

	// Test confirming non-active task
	task := db.Task{
		ID:         999,
		TemplateID: 1,
		ProjectID:  1,
	}

	err := pool.ConfirmTask(task)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "task is not active")

	// Test rejecting non-active task
	err = pool.RejectTask(task)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "task is not active")
}
