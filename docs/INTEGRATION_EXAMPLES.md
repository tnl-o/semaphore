# Integration Examples

## Example 1: Using ConfirmActionDialog for Delete Operations

```vue
<template>
  <v-btn color="error" @click="showDeleteDialog = true">
    Delete Project
  </v-btn>
  
  <ConfirmActionDialog
    v-model="showDeleteDialog"
    type="error"
    title="Delete Project"
    :message="`Are you sure you want to delete project '${project.name}'?`"
    warning-message="All templates, tasks, and history will be permanently deleted."
    :require-confirmation="true"
    confirmation-label="Type project name to confirm"
    :confirmation-placeholder="project.name"
    confirm-text="Delete"
    @confirm="deleteProject"
  />
</template>

<script>
import ConfirmActionDialog from '@/components/ConfirmActionDialog.vue';
import { api } from '@/lib/apiClient';
import { apiCache } from '@/lib/apiCache';
import toast from '@/lib/toast';

export default {
  components: {
    ConfirmActionDialog,
  },
  data() {
    return {
      showDeleteDialog: false,
      project: { name: 'My Project' },
    };
  },
  methods: {
    async deleteProject() {
      try {
        await api.delete(`/api/project/${this.project.id}`);
        // Invalidate cache
        apiCache.invalidate(/\/api\/projects/);
        toast.success('Project deleted successfully');
        this.$router.push('/');
      } catch (error) {
        toast.apiError(error);
      }
    },
  },
};
</script>
```

## Example 2: Integrating Ansible Components in Task View

```vue
<template>
  <v-container>
    <v-tabs v-model="tab">
      <v-tab>Logs</v-tab>
      <v-tab>Facts</v-tab>
      <v-tab>Structure</v-tab>
    </v-tabs>
    
    <v-tabs-items v-model="tab">
      <v-tab-item>
        <AnsibleLogViewer :logs="taskLogs" />
      </v-tab-item>
      
      <v-tab-item>
        <AnsibleFactsViewer :facts="taskFacts" />
      </v-tab-item>
      
      <v-tab-item>
        <PlaybookStructureViewer :playbook="parsedPlaybook" />
      </v-tab-item>
    </v-tabs-items>
  </v-container>
</template>

<script>
import { ansibleParser } from '@/lib/ansibleParser';
import AnsibleLogViewer from '@/components/AnsibleLogViewer.vue';
import AnsibleFactsViewer from '@/components/AnsibleFactsViewer.vue';
import PlaybookStructureViewer from '@/components/PlaybookStructureViewer.vue';
import { api } from '@/lib/apiClient';

export default {
  components: {
    AnsibleLogViewer,
    AnsibleFactsViewer,
    PlaybookStructureViewer,
  },
  data() {
    return {
      tab: 0,
      taskLogs: [],
      taskFacts: {},
      parsedPlaybook: null,
    };
  },
  async created() {
    await this.loadTaskData();
  },
  methods: {
    async loadTaskData() {
      const logs = await api.get(`/api/project/${this.projectId}/tasks/${this.taskId}/output`);
      this.taskLogs = logs.data;
      
      // Parse Ansible output
      const parsed = ansibleParser.parse(this.taskLogs);
      this.parsedPlaybook = parsed;
      this.taskFacts = parsed.facts;
    },
  },
};
</script>
```

## Example 3: Using SecurePasswordInput in Forms

```vue
<template>
  <v-form>
    <SecurePasswordInput
      v-model="password"
      label="Password"
      :rules="[v => !!v || 'Password is required']"
      required
    />
  </v-form>
</template>

<script>
import SecurePasswordInput from '@/components/SecurePasswordInput.vue';

export default {
  components: {
    SecurePasswordInput,
  },
  data() {
    return {
      password: '',
    };
  },
};
</script>
```

## Example 4: Using ContentSkeleton for Loading States

```vue
<template>
  <ContentSkeleton :loading="tasks === null" type="table">
    <TaskList :tasks="tasks" />
  </ContentSkeleton>
</template>

<script>
import ContentSkeleton from '@/components/ContentSkeleton.vue';
import TaskList from '@/components/TaskList.vue';
import { api } from '@/lib/apiClient';

export default {
  components: {
    ContentSkeleton,
    TaskList,
  },
  data() {
    return {
      tasks: null,
    };
  },
  async created() {
    this.tasks = await api.get('/api/tasks').then(r => r.data);
  },
};
</script>
```

## Example 5: Using OptimisticUpdate Mixin

```vue
<script>
import OptimisticUpdate from '@/mixins/OptimisticUpdate';
import { api } from '@/lib/apiClient';
import { apiCache } from '@/lib/apiCache';

export default {
  mixins: [OptimisticUpdate],
  methods: {
    async stopTask(taskId) {
      await this.optimisticUpdate(
        () => {
          // Optimistic update
          const task = this.tasks.find(t => t.id === taskId);
          const previousStatus = task.status;
          task.status = 'stopping';
          return { task, previousStatus };
        },
        async () => {
          // API call
          await api.post(`/api/project/${this.projectId}/tasks/${taskId}/stop`);
          // Invalidate cache
          apiCache.invalidate(/\/api\/project\/\d+\/tasks/);
        },
        (previousState) => {
          // Rollback on error
          previousState.task.status = previousState.previousStatus;
        }
      );
    },
  },
};
</script>
```

