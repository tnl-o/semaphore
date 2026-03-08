# New Features Implementation Plan

This document outlines plans for implementing new features in Semaphore UI.

## Overview

This plan covers four major new features:
1. Export/Import configurations
2. Dashboard with metrics
3. Enhanced notification system
4. Task rollback functionality

## 1. Export/Import Configurations

### Goal

Allow users to export and import project configurations for backup, migration, or sharing.

### Requirements

- Export project configuration to JSON/YAML
- Import project configuration from JSON/YAML
- Validate imported configuration
- Handle dependencies (templates, inventories, etc.)
- Support partial imports

### Implementation Plan

#### 1.1 Export Functionality

**Backend API:**

```go
// api/projects/export.go
func ExportProject(w http.ResponseWriter, r *http.Request) {
    project := helpers.GetFromContext(r, "project").(*db.Project)
    
    export := ProjectExport{
        Project: project,
        Templates: getTemplates(project.ID),
        Inventories: getInventories(project.ID),
        Repositories: getRepositories(project.ID),
        Keys: getKeys(project.ID),
        // ... other resources
    }
    
    format := r.URL.Query().Get("format") // json or yaml
    if format == "yaml" {
        w.Header().Set("Content-Type", "application/x-yaml")
        yaml.NewEncoder(w).Encode(export)
    } else {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(export)
    }
}
```

**Frontend:**

```vue
<template>
  <v-btn @click="exportProject">
    <v-icon left>mdi-download</v-icon>
    Export Project
  </v-btn>
</template>

<script>
export default {
  methods: {
    async exportProject() {
      const response = await api.get(`/api/project/${this.projectId}/export?format=json`);
      const blob = new Blob([JSON.stringify(response.data, null, 2)], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `project-${this.projectId}.json`;
      a.click();
    }
  }
}
</script>
```

#### 1.2 Import Functionality

**Backend API:**

```go
// api/projects/import.go
func ImportProject(w http.ResponseWriter, r *http.Request) {
    var importData ProjectImport
    if !helpers.Bind(w, r, &importData) {
        return
    }
    
    // Validate import data
    if err := validateImport(importData); err != nil {
        helpers.WriteError(w, err)
        return
    }
    
    // Import project
    project, err := importProjectData(importData)
    if err != nil {
        helpers.WriteError(w, err)
        return
    }
    
    helpers.WriteJSON(w, http.StatusCreated, project)
}
```

**Frontend:**

```vue
<template>
  <v-file-input
    v-model="file"
    label="Import Project"
    accept=".json,.yaml,.yml"
    @change="importProject"
  ></v-file-input>
</template>
```

### Checklist

- [ ] Design export/import data structure
- [ ] Implement export API endpoint
- [ ] Implement import API endpoint
- [ ] Add validation for imported data
- [ ] Create export UI component
- [ ] Create import UI component
- [ ] Add error handling
- [ ] Add tests
- [ ] Document format specification

## 2. Dashboard with Metrics

### Goal

Create a dashboard showing project statistics, task execution metrics, and system health.

### Requirements

- Project statistics (total tasks, success rate, etc.)
- Task execution trends
- System health indicators
- Customizable widgets
- Real-time updates

### Implementation Plan

#### 2.1 Backend Metrics API

```go
// api/dashboard.go
func GetDashboardMetrics(w http.ResponseWriter, r *http.Request) {
    project := helpers.GetFromContext(r, "project").(*db.Project)
    user := helpers.GetFromContext(r, "user").(*db.User)
    
    metrics := DashboardMetrics{
        ProjectStats: getProjectStats(project.ID),
        TaskStats: getTaskStats(project.ID),
        RecentTasks: getRecentTasks(project.ID, 10),
        SystemHealth: getSystemHealth(),
    }
    
    helpers.WriteJSON(w, http.StatusOK, metrics)
}
```

#### 2.2 Frontend Dashboard

```vue
<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="4">
        <StatCard
          title="Total Tasks"
          :value="metrics.totalTasks"
          icon="mdi-play-circle"
        />
      </v-col>
      <v-col cols="12" md="4">
        <StatCard
          title="Success Rate"
          :value="metrics.successRate"
          suffix="%"
          icon="mdi-check-circle"
        />
      </v-col>
      <v-col cols="12" md="4">
        <StatCard
          title="Active Tasks"
          :value="metrics.activeTasks"
          icon="mdi-clock-outline"
        />
      </v-col>
    </v-row>
    
    <v-row>
      <v-col cols="12">
        <TaskTrendChart :data="metrics.taskTrends" />
      </v-col>
    </v-row>
  </v-container>
</template>
```

### Checklist

- [ ] Design metrics data structure
- [ ] Implement metrics API endpoints
- [ ] Create dashboard view
- [ ] Create stat card component
- [ ] Create chart components
- [ ] Add real-time updates
- [ ] Add widget customization
- [ ] Add tests
- [ ] Document metrics

## 3. Enhanced Notification System

### Goal

Expand notification channels beyond email to include webhooks, Slack, Teams, and more.

### Requirements

- Multiple notification channels
- Configurable notification rules
- Template support
- Delivery status tracking
- Retry mechanism

### Implementation Plan

#### 3.1 Notification Service

```go
// services/notifications/notification.go
type NotificationService interface {
    Send(notification Notification) error
    SendAsync(notification Notification)
    RegisterChannel(channel NotificationChannel)
}

type NotificationChannel interface {
    Send(notification Notification) error
    Validate(config map[string]interface{}) error
}
```

#### 3.2 Webhook Channel

```go
// services/notifications/webhook.go
type WebhookChannel struct {
    client *http.Client
}

func (w *WebhookChannel) Send(notification Notification) error {
    req, err := http.NewRequest("POST", notification.Config.URL, bytes.NewBuffer(notification.Payload))
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    if notification.Config.Secret != "" {
        req.Header.Set("X-Webhook-Secret", notification.Config.Secret)
    }
    
    resp, err := w.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        return fmt.Errorf("webhook returned status %d", resp.StatusCode)
    }
    
    return nil
}
```

#### 3.3 Slack Channel

```go
// services/notifications/slack.go
type SlackChannel struct {
    client *http.Client
}

func (s *SlackChannel) Send(notification Notification) error {
    payload := map[string]interface{}{
        "text": notification.Message,
        "channel": notification.Config.Channel,
    }
    
    // Send to Slack webhook
    // Implementation...
}
```

### Checklist

- [ ] Design notification service interface
- [ ] Implement webhook channel
- [ ] Implement Slack channel
- [ ] Implement Teams channel
- [ ] Add notification rules
- [ ] Add template support
- [ ] Add delivery tracking
- [ ] Create UI for configuration
- [ ] Add tests
- [ ] Document channels

## 4. Task Rollback Functionality

### Goal

Allow users to rollback completed tasks to previous states.

### Requirements

- Track task execution history
- Store task state snapshots
- Rollback to previous state
- Validate rollback feasibility
- Audit trail

### Implementation Plan

#### 4.1 Task History

```go
// db/TaskHistory.go
type TaskHistory struct {
    ID        int       `db:"id" json:"id"`
    TaskID    int       `db:"task_id" json:"task_id"`
    State     string    `db:"state" json:"state"`
    Snapshot  []byte    `db:"snapshot" json:"-"` // JSON snapshot
    Created   time.Time `db:"created" json:"created"`
    CreatedBy int       `db:"created_by" json:"created_by"`
}
```

#### 4.2 Rollback API

```go
// api/tasks/rollback.go
func RollbackTask(w http.ResponseWriter, r *http.Request) {
    task := helpers.GetFromContext(r, "task").(*db.Task)
    historyID, _ := helpers.GetIntParam("history_id", w, r)
    
    // Get history snapshot
    history, err := store.GetTaskHistory(historyID)
    if err != nil {
        helpers.WriteError(w, err)
        return
    }
    
    // Validate rollback
    if err := validateRollback(task, history); err != nil {
        helpers.WriteError(w, err)
        return
    }
    
    // Perform rollback
    if err := performRollback(task, history); err != nil {
        helpers.WriteError(w, err)
        return
    }
    
    helpers.WriteJSON(w, http.StatusOK, map[string]string{
        "message": "Task rolled back successfully",
    })
}
```

### Checklist

- [ ] Design task history schema
- [ ] Implement task snapshot storage
- [ ] Implement rollback API
- [ ] Add rollback validation
- [ ] Create rollback UI
- [ ] Add audit logging
- [ ] Add tests
- [ ] Document rollback process

## Implementation Priority

1. **High Priority:**
   - Export/Import configurations
   - Dashboard with metrics

2. **Medium Priority:**
   - Enhanced notification system

3. **Low Priority:**
   - Task rollback functionality

## Resources

- [JSON Schema](https://json-schema.org/)
- [YAML Specification](https://yaml.org/spec/)
- [Chart.js Documentation](https://www.chartjs.org/)
- [Webhook Best Practices](https://webhooks.fyi/)

