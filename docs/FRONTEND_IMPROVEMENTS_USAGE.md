# Frontend Improvements - Usage Guide

## Quick Start

### 1. Install Dependencies

```bash
cd web
npm install
```

This will install `dompurify` which is required for XSS protection.

### 2. Using New Components

#### SafeLogOutput (XSS Protection)

Replace all instances of `v-html` with log output:

**Before:**
```vue
<div v-html="$options.filters.formatLog(record.output)"></div>
```

**After:**
```vue
<SafeLogOutput :output="record.output" />
```

#### API Client (Centralized)

Replace direct `axios` usage:

**Before:**
```javascript
import axios from 'axios';

const response = await axios.get('/api/project/1/templates');
```

**After:**
```javascript
import { api } from '@/lib/apiClient';

const response = await api.get('/api/project/1/templates');
// Response is automatically cached for GET requests
```

#### Confirmation Dialogs

Add confirmation for destructive actions:

```vue
<template>
  <v-btn @click="showDeleteDialog = true">Delete</v-btn>
  
  <ConfirmActionDialog
    v-model="showDeleteDialog"
    type="error"
    title="Delete Project"
    :message="`Are you sure you want to delete '${project.name}'?`"
    warning-message="This action cannot be undone."
    :require-confirmation="true"
    confirmation-label="Type project name to confirm"
    :confirmation-placeholder="project.name"
    confirm-text="Delete"
    @confirm="deleteProject"
  />
</template>

<script>
export default {
  data() {
    return {
      showDeleteDialog: false,
    };
  },
  methods: {
    async deleteProject() {
      // Delete logic here
    },
  },
};
</script>
```

#### Breadcrumbs

Already integrated in `App.vue`. Automatically displays based on current route.

#### Ansible Components

**Ansible Facts Viewer:**
```vue
<AnsibleFactsViewer :facts="taskFacts" />
```

**Ansible Tags Manager:**
```vue
<AnsibleTagsManager
  :playbook="parsedPlaybook"
  v-model="selectedTags"
  :skip-tags.sync="skipTags"
/>
```

**Ansible Log Viewer:**
```vue
<AnsibleLogViewer :logs="taskLogs" />
```

**Ansible Inventory Viewer:**
```vue
<AnsibleInventoryViewer :inventory="inventoryData" />
```

### 3. Cache Invalidation

After mutations, invalidate relevant cache:

```javascript
import { apiCache } from '@/lib/apiCache';
import { api } from '@/lib/apiClient';

// After creating/updating/deleting a template
await api.post('/api/project/1/templates', data);
apiCache.invalidate(/\/api\/project\/\d+\/templates/);

// After task status change
apiCache.invalidate(/\/api\/project\/\d+\/tasks/);
```

### 4. Secret Masking

When logging sensitive data:

```javascript
import { maskSecrets } from '@/lib/secretMasker';

console.log(maskSecrets(logOutput));
// Secrets will be replaced with ***MASKED***
```

### 5. Performance Monitoring

Measure function performance:

```javascript
import { measurePerformance } from '@/lib/performance';

const data = await measurePerformance('loadTemplates', () =>
  api.get('/api/project/1/templates')
);
```

### 6. Optimistic UI Updates

For immediate UI feedback:

```vue
<script>
import OptimisticUpdate from '@/mixins/OptimisticUpdate';
import { api } from '@/lib/apiClient';

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
        () => api.post(`/tasks/${taskId}/stop`),
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

## Migration Checklist

- [ ] Replace all `axios` imports with `@/lib/apiClient`
- [ ] Replace `v-html` log rendering with `SafeLogOutput`
- [ ] Add `ConfirmActionDialog` to delete operations
- [ ] Update password inputs to use `SecurePasswordInput`
- [ ] Add cache invalidation after mutations
- [ ] Integrate `ContentSkeleton` in async components
- [ ] Add `AnsibleFactsViewer` to task details
- [ ] Use `AnsibleLogViewer` as alternative log view

## Component Reference

See `docs/FRONTEND_DEEP_REVIEW.md` for detailed component documentation and examples.

