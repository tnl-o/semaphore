# Deep Review: Semaphore UI Frontend Improvements

## Executive Summary

This document provides a comprehensive technical review and prioritized improvement plan for the Semaphore UI frontend. The analysis is based on the current Vue 2.6.14 + Vuetify 2.6.10 codebase and focuses on security, UX, performance, and maintainability improvements critical for DevOps teams.

**Current Stack:**
- Vue 2.6.14
- Vuetify 2.6.10
- Vue Router 3.5.4
- Axios 1.12.0
- Vue CLI (Webpack)
- No TypeScript
- Minimal test coverage

---

## 🚨 P0: Critical Security & Reliability Issues

### P0.1: XSS Vulnerability in Log Rendering

**Issue:** Logs are rendered using `v-html` without sanitization, creating XSS attack vectors.

**Current Code:**
```vue
<!-- TaskLogView.vue, TaskLogViewRecord.vue -->
<div class="task-log-records__output" v-html="$options.filters.formatLog(record.output)">
```

**Risk:** Ansible playbook output or malicious task names could inject JavaScript, compromising user sessions or stealing credentials.

**Solution:**
1. Sanitize HTML output using DOMPurify
2. Create a safe log renderer component
3. Add Content Security Policy headers

**Implementation:**

```javascript
// web/src/lib/sanitize.js
import DOMPurify from 'dompurify';

/**
 * Sanitize HTML content for safe rendering
 * @param {string} html - HTML string to sanitize
 * @returns {string} - Sanitized HTML
 */
export function sanitizeHtml(html) {
  if (!html) return '';
  
  return DOMPurify.sanitize(html, {
    ALLOWED_TAGS: ['span', 'div', 'br', 'pre'],
    ALLOWED_ATTR: ['class', 'style'],
    ALLOW_DATA_ATTR: false,
    // Preserve ANSI color classes
    ALLOWED_CLASSES: {
      'span': /^ansi-(black|red|green|yellow|blue|magenta|cyan|white|bright-.*)$/,
    },
  });
}
```

```vue
<!-- web/src/components/SafeLogOutput.vue -->
<template>
  <div 
    class="task-log-records__output" 
    v-html="sanitizedOutput"
  />
</template>

<script>
import { sanitizeHtml } from '@/lib/sanitize';
import { AnsiUp } from 'ansi_up';

const ansiConverter = new AnsiUp();
// ... configure ansiConverter colors ...

export default {
  props: {
    output: {
      type: String,
      required: true,
    },
  },
  computed: {
    sanitizedOutput() {
      const html = ansiConverter.ansi_to_html(String(this.output || ''));
      return sanitizeHtml(html);
    },
  },
};
</script>
```

**Quick Win:** ✅ Can be implemented in < 4 hours

---

### P0.2: Centralized API Client with Security Interceptors

**Issue:** Axios is used directly throughout components without centralized error handling, retry logic, or request/response interceptors.

**Current State:**
- No automatic token refresh
- No request retry on network failures
- Inconsistent error handling
- No request/response logging for debugging

**Solution:** Create a typed API client with interceptors.

**Implementation:**

```javascript
// web/src/lib/apiClient.js
import axios from 'axios';
import router from '@/router';
import { getErrorMessage } from '@/lib/error';
import toast from '@/lib/toast';

// Create axios instance
const apiClient = axios.create({
  baseURL: document.baseURI,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - add auth token, logging
apiClient.interceptors.request.use(
  (config) => {
    // Add correlation ID for request tracking
    config.headers['X-Request-ID'] = generateRequestId();
    
    // Log request in development
    if (process.env.NODE_ENV === 'development') {
      console.log(`[API] ${config.method.toUpperCase()} ${config.url}`, config.data);
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor - handle errors, retries
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    // Handle 401 - redirect to login
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      // Clear auth state
      // Redirect to login with return URL
      router.push({
        path: '/auth/login',
        query: { redirect: router.currentRoute.fullPath },
      });
      
      return Promise.reject(error);
    }
    
    // Handle 429 - rate limiting with exponential backoff
    if (error.response?.status === 429 && !originalRequest._retry) {
      originalRequest._retry = true;
      const retryAfter = error.response.headers['retry-after'] || 1;
      
      await new Promise((resolve) => setTimeout(resolve, retryAfter * 1000));
      return apiClient(originalRequest);
    }
    
    // Handle network errors with retry
    if (!error.response && !originalRequest._retry && originalRequest.method === 'get') {
      originalRequest._retry = true;
      await new Promise((resolve) => setTimeout(resolve, 1000));
      return apiClient(originalRequest);
    }
    
    // Log error
    console.error('[API Error]', {
      url: originalRequest.url,
      method: originalRequest.method,
      status: error.response?.status,
      message: getErrorMessage(error),
    });
    
    return Promise.reject(error);
  }
);

// Helper to generate request ID
function generateRequestId() {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

// Export typed methods
export const api = {
  get: (url, config) => apiClient.get(url, config),
  post: (url, data, config) => apiClient.post(url, data, config),
  put: (url, data, config) => apiClient.put(url, data, config),
  delete: (url, config) => apiClient.delete(url, config),
  patch: (url, data, config) => apiClient.patch(url, data, config),
};

export default apiClient;
```

**Migration Path:**
1. Replace `axios` imports with `@/lib/apiClient`
2. Update all API calls to use the new client
3. Remove manual error handling (now centralized)

**Quick Win:** ✅ Can be implemented in < 6 hours

---

### P0.3: Secret Masking in UI

**Issue:** Secrets (passwords, SSH keys, vault tokens) may be exposed in:
- Form inputs (even with `type="password"`, values are in DOM)
- Browser DevTools
- Task logs (if accidentally logged)
- Error messages

**Current Risk Areas:**
- `SecretStorageForm.vue` - vault tokens
- `KeyForm.vue` - SSH private keys
- `EnvironmentForm.vue` - secret environment variables
- `TaskForm.vue` - survey vars with `type="secret"`

**Solution:**
1. Never store secrets in component data
2. Use secure input components with memory clearing
3. Mask secrets in logs and error messages
4. Add Content Security Policy

**Implementation:**

```vue
<!-- web/src/components/SecurePasswordInput.vue -->
<template>
  <v-text-field
    :value="displayValue"
    @input="handleInput"
    :type="showPassword ? 'text' : 'password'"
    :append-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
    @click:append="showPassword = !showPassword"
    v-bind="$attrs"
    v-on="$listeners"
  />
</template>

<script>
export default {
  name: 'SecurePasswordInput',
  inheritAttrs: false,
  props: {
    value: String,
  },
  data() {
    return {
      showPassword: false,
      // Never store actual value in component data
      internalValue: '',
    };
  },
  computed: {
    displayValue() {
      // Show masked value when not focused
      if (!this.showPassword && this.value) {
        return '•'.repeat(Math.min(this.value.length, 20));
      }
      return this.value || '';
    },
  },
  methods: {
    handleInput(value) {
      // Immediately emit, don't store
      this.$emit('input', value);
      // Clear internal reference after next tick
      this.$nextTick(() => {
        this.internalValue = '';
      });
    },
  },
  beforeDestroy() {
    // Clear any references
    this.internalValue = '';
  },
};
</script>
```

**Additional Measures:**
```javascript
// web/src/lib/secretMasker.js
/**
 * Mask secrets in strings (for logging/debugging)
 */
export function maskSecrets(text, patterns = []) {
  if (!text) return text;
  
  const defaultPatterns = [
    /password["\s:=]+([^\s"']+)/gi,
    /token["\s:=]+([^\s"']+)/gi,
    /secret["\s:=]+([^\s"']+)/gi,
    /api[_-]?key["\s:=]+([^\s"']+)/gi,
    /-----BEGIN.*?-----[\s\S]*?-----END.*?-----/gi, // SSH keys
  ];
  
  let masked = text;
  [...defaultPatterns, ...patterns].forEach((pattern) => {
    masked = masked.replace(pattern, (match, secret) => {
      if (secret) {
        return match.replace(secret, '***MASKED***');
      }
      return '***MASKED***';
    });
  });
  
  return masked;
}
```

**Quick Win:** ⚠️ Moderate effort (8-12 hours) - requires component updates

---

## ⚠️ P1: High-Priority UX & Architecture Improvements

### P1.1: TypeScript Migration (Gradual)

**Why:** Type safety prevents bugs, improves IDE support, and makes refactoring safer.

**Current State:** 100% JavaScript, no type checking

**Strategy:** Gradual migration using `// @ts-check` and `.ts` files

**Phase 1: Setup (Quick Win - 4 hours)**

```json
// web/tsconfig.json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "moduleResolution": "node",
    "strict": true,
    "jsx": "preserve",
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    },
    "types": ["node", "jest"]
  },
  "include": [
    "src/**/*.ts",
    "src/**/*.tsx",
    "src/**/*.vue"
  ],
  "exclude": [
    "node_modules",
    "dist"
  ]
}
```

```javascript
// web/vue.config.js - add TypeScript support
module.exports = {
  // ... existing config
  chainWebpack: (config) => {
    // ... existing chainWebpack
    config.resolve.extensions
      .merge(['.ts', '.tsx']);
    
    config.module
      .rule('ts')
      .test(/\.tsx?$/)
      .use('ts-loader')
      .loader('ts-loader')
      .options({
        appendTsSuffixTo: [/\.vue$/],
        transpileOnly: true,
      });
  },
};
```

**Phase 2: Migrate Utilities First (16 hours)**

Start with type-safe utilities:
- `lib/error.js` → `lib/error.ts`
- `lib/apiClient.js` → `lib/apiClient.ts`
- `lib/toast.js` → `lib/toast.ts`

**Phase 3: API Types from OpenAPI (8 hours)**

Generate TypeScript types from OpenAPI spec:

```bash
# Install openapi-typescript
npm install -D openapi-typescript

# Generate types
npx openapi-typescript api-docs.yml -o src/types/api.d.ts
```

**Phase 4: Component Migration (40+ hours)**

Migrate components gradually, starting with leaf components.

**Quick Win:** ✅ Setup can be done in < 4 hours, provides immediate value

---

### P1.2: Enhanced Log Viewer with Collapsible Tasks

**Current State:** Logs are displayed as a flat list with ANSI colors. No grouping by Ansible tasks.

**Improvement:** Group logs by Ansible play/task with collapsible sections.

**Implementation:**

```vue
<!-- web/src/components/AnsibleLogViewer.vue -->
<template>
  <div class="ansible-log-viewer">
    <div
      v-for="(task, index) in parsedTasks"
      :key="task.id"
      class="ansible-task"
    >
      <div
        class="ansible-task__header"
        @click="toggleTask(task.id)"
        :class="{ 'ansible-task__header--expanded': task.expanded }"
      >
        <v-icon small class="mr-2">
          {{ task.expanded ? 'mdi-chevron-down' : 'mdi-chevron-right' }}
        </v-icon>
        <TaskStatusBadge :status="task.status" />
        <span class="ml-2">{{ task.name }}</span>
        <v-spacer />
        <span class="text-caption text--secondary">
          {{ task.duration }}
        </span>
      </div>
      
      <v-expand-transition>
        <div v-show="task.expanded" class="ansible-task__content">
          <div
            v-for="record in task.records"
            :key="record.id"
            class="log-record"
          >
            <span class="log-record__time">{{ record.time | formatTime }}</span>
            <SafeLogOutput :output="record.output" />
          </div>
        </div>
      </v-expand-transition>
    </div>
  </div>
</template>

<script>
import SafeLogOutput from './SafeLogOutput.vue';
import TaskStatusBadge from './TaskStatusBadge.vue';

export default {
  components: {
    SafeLogOutput,
    TaskStatusBadge,
  },
  props: {
    logs: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      expandedTasks: new Set(),
    };
  },
  computed: {
    parsedTasks() {
      // Parse Ansible output to extract task boundaries
      // Look for patterns like:
      // "TASK [task name] ***"
      // "PLAY [play name] ***"
      const tasks = [];
      let currentTask = null;
      
      this.logs.forEach((log) => {
        const taskMatch = log.output.match(/^(TASK|PLAY)\s+\[([^\]]+)\]/);
        
        if (taskMatch) {
          // Save previous task
          if (currentTask) {
            tasks.push(currentTask);
          }
          
          // Start new task
          currentTask = {
            id: `${taskMatch[1]}-${taskMatch[2]}-${log.time}`,
            type: taskMatch[1],
            name: taskMatch[2],
            records: [log],
            status: 'running',
            expanded: this.expandedTasks.has(`${taskMatch[1]}-${taskMatch[2]}-${log.time}`),
          };
        } else if (currentTask) {
          currentTask.records.push(log);
          
          // Detect task completion
          if (log.output.match(/^(ok|changed|failed|skipped):/)) {
            const statusMatch = log.output.match(/^(ok|changed|failed|skipped):/);
            if (statusMatch) {
              currentTask.status = statusMatch[1] === 'failed' ? 'error' : 'success';
            }
          }
        } else {
          // Logs before first task
          if (!tasks.length || tasks[tasks.length - 1].name !== '_prelude') {
            tasks.push({
              id: '_prelude',
              name: 'Prelude',
              records: [],
              status: 'success',
              expanded: false,
            });
          }
          tasks[tasks.length - 1].records.push(log);
        }
      });
      
      if (currentTask) {
        tasks.push(currentTask);
      }
      
      return tasks;
    },
  },
  methods: {
    toggleTask(taskId) {
      if (this.expandedTasks.has(taskId)) {
        this.expandedTasks.delete(taskId);
      } else {
        this.expandedTasks.add(taskId);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.ansible-task {
  border-left: 3px solid #ccc;
  margin-bottom: 8px;
  
  &__header {
    padding: 8px 12px;
    cursor: pointer;
    display: flex;
    align-items: center;
    background: rgba(0, 0, 0, 0.05);
    transition: background 0.2s;
    
    &:hover {
      background: rgba(0, 0, 0, 0.1);
    }
    
    &--expanded {
      border-left-color: #1976d2;
    }
  }
  
  &__content {
    padding-left: 32px;
    background: #000;
    color: #fff;
    font-family: monospace;
  }
}

.log-record {
  display: flex;
  padding: 2px 8px;
  
  &__time {
    width: 120px;
    color: #888;
    flex-shrink: 0;
  }
}
</style>
```

**Quick Win:** ✅ Can be implemented in < 8 hours

---

### P1.3: Breadcrumb Navigation

**Current State:** No breadcrumbs, users rely on browser back button or sidebar navigation.

**Implementation:**

```vue
<!-- web/src/components/Breadcrumbs.vue -->
<template>
  <v-breadcrumbs
    :items="breadcrumbItems"
    class="pa-0"
    large
  >
    <template v-slot:divider>
      <v-icon>mdi-chevron-right</v-icon>
    </template>
    
    <template v-slot:item="{ item }">
      <v-breadcrumbs-item
        :to="item.to"
        :disabled="item.disabled"
        :exact="item.exact"
      >
        <v-icon v-if="item.icon" small class="mr-1">{{ item.icon }}</v-icon>
        {{ item.text }}
      </v-breadcrumbs-item>
    </template>
  </v-breadcrumbs>
</template>

<script>
export default {
  name: 'Breadcrumbs',
  computed: {
    breadcrumbItems() {
      const route = this.$route;
      const items = [];
      
      // Home
      items.push({
        text: this.$t('dashboard'),
        to: '/',
        icon: 'mdi-home',
        exact: true,
      });
      
      // Project level
      if (route.params.projectId) {
        const project = this.$store?.state?.projects?.find(
          p => p.id === parseInt(route.params.projectId)
        ) || { name: `Project ${route.params.projectId}` };
        
        items.push({
          text: project.name,
          to: `/project/${route.params.projectId}`,
          icon: 'mdi-folder',
        });
        
        // Template level
        if (route.params.templateId) {
          const template = this.$store?.state?.templates?.find(
            t => t.id === parseInt(route.params.templateId)
          ) || { name: `Template ${route.params.templateId}` };
          
          items.push({
            text: template.name,
            to: route.path,
            icon: 'mdi-file-document',
            disabled: true,
          });
        }
        
        // Other project sections
        const sectionMap = {
          'templates': { text: this.$t('templates'), icon: 'mdi-file-document-outline' },
          'history': { text: this.$t('history'), icon: 'mdi-history' },
          'settings': { text: this.$t('settings'), icon: 'mdi-cog' },
          'team': { text: this.$t('team'), icon: 'mdi-account-group' },
          'keys': { text: this.$t('keys'), icon: 'mdi-key' },
          'inventory': { text: this.$t('inventory'), icon: 'mdi-database' },
          'environment': { text: this.$t('environment'), icon: 'mdi-code-braces' },
        };
        
        const section = Object.keys(sectionMap).find(
          s => route.path.includes(`/${s}`)
        );
        
        if (section && !route.params.templateId) {
          items.push({
            ...sectionMap[section],
            to: route.path,
            disabled: true,
          });
        }
      }
      
      return items;
    },
  },
};
</script>
```

**Usage in App.vue:**
```vue
<v-app-bar>
  <Breadcrumbs />
  <v-spacer />
  <!-- ... other app bar content -->
</v-app-bar>
```

**Quick Win:** ✅ Can be implemented in < 4 hours

---

### P1.4: Confirmation Dialogs for Critical Actions

**Current State:** No confirmations for:
- Deleting projects
- Stopping running tasks
- Deleting templates
- Removing team members

**Implementation:**

```vue
<!-- web/src/components/ConfirmActionDialog.vue -->
<template>
  <v-dialog
    :value="value"
    @input="$emit('input', $event)"
    max-width="500"
    persistent
  >
    <v-card>
      <v-card-title class="headline">
        <v-icon :color="type" class="mr-2">{{ icon }}</v-icon>
        {{ title }}
      </v-card-title>
      
      <v-card-text>
        <p>{{ message }}</p>
        
        <!-- Additional warning for destructive actions -->
        <v-alert
          v-if="type === 'error'"
          type="warning"
          dense
          outlined
          class="mt-3"
        >
          {{ warningMessage || 'This action cannot be undone.' }}
        </v-alert>
        
        <!-- Confirmation input for very destructive actions -->
        <v-text-field
          v-if="requireConfirmation"
          v-model="confirmationText"
          :label="confirmationLabel"
          :placeholder="confirmationPlaceholder"
          outlined
          dense
          class="mt-3"
        />
      </v-card-text>
      
      <v-card-actions>
        <v-spacer />
        <v-btn
          text
          @click="$emit('input', false)"
        >
          Cancel
        </v-btn>
        <v-btn
          :color="type"
          :disabled="requireConfirmation && confirmationText !== confirmationPlaceholder"
          @click="confirm"
        >
          {{ confirmText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: 'ConfirmActionDialog',
  props: {
    value: Boolean,
    title: {
      type: String,
      required: true,
    },
    message: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      default: 'primary',
      validator: (v) => ['primary', 'warning', 'error'].includes(v),
    },
    confirmText: {
      type: String,
      default: 'Confirm',
    },
    warningMessage: String,
    requireConfirmation: Boolean,
    confirmationLabel: String,
    confirmationPlaceholder: String,
  },
  data() {
    return {
      confirmationText: '',
    };
  },
  computed: {
    icon() {
      const icons = {
        primary: 'mdi-information',
        warning: 'mdi-alert',
        error: 'mdi-alert-circle',
      };
      return icons[this.type];
    },
  },
  watch: {
    value(newVal) {
      if (!newVal) {
        // Reset confirmation on close
        this.confirmationText = '';
      }
    },
  },
  methods: {
    confirm() {
      this.$emit('confirm');
      this.$emit('input', false);
    },
  },
};
</script>
```

**Usage Example:**
```vue
<template>
  <div>
    <v-btn @click="showDeleteDialog = true">Delete Project</v-btn>
    
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
  </div>
</template>
```

**Quick Win:** ✅ Can be implemented in < 4 hours

---

### P1.5: Virtualized Task Lists

**Current State:** `TaskList.vue` loads all tasks into a `v-data-table` without virtualization.

**Issue:** Performance degrades with 100+ tasks.

**Solution:** Use `vue-virtual-scroll-list` (already in dependencies) or `vue-virtual-scroller`.

**Implementation:**

```vue
<!-- web/src/components/VirtualizedTaskList.vue -->
<template>
  <div class="virtualized-task-list">
    <v-toolbar flat dense>
      <v-text-field
        v-model="searchQuery"
        prepend-inner-icon="mdi-magnify"
        label="Search tasks"
        single-line
        hide-details
        clearable
        class="mr-4"
      />
      <v-select
        v-model="statusFilter"
        :items="statusOptions"
        label="Status"
        multiple
        clearable
        dense
        hide-details
        class="mr-4"
        style="max-width: 200px;"
      />
    </v-toolbar>
    
    <VirtualList
      :data-key="'id'"
      :data-sources="filteredTasks"
      :data-component="TaskRow"
      :estimate-size="60"
      :keeps="50"
      style="height: calc(100vh - 200px);"
    />
  </div>
</template>

<script>
import VirtualList from 'vue-virtual-scroll-list';
import TaskRow from './TaskRow.vue';

export default {
  components: {
    VirtualList,
  },
  props: {
    tasks: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      TaskRow,
      searchQuery: '',
      statusFilter: [],
      statusOptions: [
        { text: 'Success', value: 'success' },
        { text: 'Error', value: 'error' },
        { text: 'Running', value: 'running' },
        { text: 'Stopped', value: 'stopped' },
      ],
    };
  },
  computed: {
    filteredTasks() {
      let filtered = this.tasks;
      
      // Search filter
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase();
        filtered = filtered.filter((task) => {
          return (
            task.id.toString().includes(query) ||
            (task.message || '').toLowerCase().includes(query) ||
            (task.commit_message || '').toLowerCase().includes(query)
          );
        });
      }
      
      // Status filter
      if (this.statusFilter.length > 0) {
        filtered = filtered.filter((task) =>
          this.statusFilter.includes(task.status)
        );
      }
      
      return filtered;
    },
  },
};
</script>
```

**Quick Win:** ✅ Can be implemented in < 6 hours (component already exists)

---

## 📋 P2: Important Improvements

### P2.1: WCAG AA Accessibility Compliance

**Current Gaps:**
- Missing ARIA labels on interactive elements
- Insufficient color contrast
- No keyboard navigation hints
- Tables without proper roles

**Implementation Checklist:**

1. **ARIA Labels:**
```vue
<v-btn
  :aria-label="$t('deleteProject')"
  icon
  @click="deleteProject"
>
  <v-icon>mdi-delete</v-icon>
</v-btn>
```

2. **Keyboard Navigation:**
```vue
<v-data-table
  :items="tasks"
  :headers="headers"
  role="table"
  aria-label="Task list"
  @keydown.native="handleKeyboardNavigation"
>
```

3. **Focus Management:**
```javascript
// Focus trap for modals
export function trapFocus(element) {
  const focusableElements = element.querySelectorAll(
    'a[href], button:not([disabled]), textarea, input, select'
  );
  const firstElement = focusableElements[0];
  const lastElement = focusableElements[focusableElements.length - 1];
  
  element.addEventListener('keydown', (e) => {
    if (e.key === 'Tab') {
      if (e.shiftKey && document.activeElement === firstElement) {
        lastElement.focus();
        e.preventDefault();
      } else if (!e.shiftKey && document.activeElement === lastElement) {
        firstElement.focus();
        e.preventDefault();
      }
    }
  });
  
  firstElement.focus();
}
```

**Quick Win:** ⚠️ Moderate effort (16-24 hours) - requires systematic audit

---

### P2.2: API Response Caching with Invalidation

**Current State:** No caching, every navigation triggers new API calls.

**Solution:** Implement request caching with smart invalidation.

**Implementation:**

```javascript
// web/src/lib/apiCache.js
class ApiCache {
  constructor(ttl = 5 * 60 * 1000) { // 5 minutes default
    this.cache = new Map();
    this.ttl = ttl;
  }
  
  get(key) {
    const entry = this.cache.get(key);
    if (!entry) return null;
    
    if (Date.now() - entry.timestamp > this.ttl) {
      this.cache.delete(key);
      return null;
    }
    
    return entry.data;
  }
  
  set(key, data) {
    this.cache.set(key, {
      data,
      timestamp: Date.now(),
    });
  }
  
  invalidate(pattern) {
    if (typeof pattern === 'string') {
      // Exact match
      this.cache.delete(pattern);
    } else if (pattern instanceof RegExp) {
      // Pattern match
      for (const key of this.cache.keys()) {
        if (pattern.test(key)) {
          this.cache.delete(key);
        }
      }
    }
  }
  
  clear() {
    this.cache.clear();
  }
}

export const apiCache = new ApiCache();

// Integrate with apiClient
import { apiCache } from '@/lib/apiCache';

apiClient.interceptors.request.use((config) => {
  // Only cache GET requests
  if (config.method === 'get' && !config.skipCache) {
    const cached = apiCache.get(config.url);
    if (cached) {
      return Promise.reject({
        ...new Error('Cached'),
        cached: true,
        data: cached,
      });
    }
  }
  return config;
});

apiClient.interceptors.response.use((response) => {
  // Cache successful GET responses
  if (response.config.method === 'get' && !response.config.skipCache) {
    apiCache.set(response.config.url, response.data);
  }
  return response;
}, (error) => {
  // Handle cache hits
  if (error.cached) {
    return Promise.resolve({
      data: error.data,
      status: 200,
      statusText: 'OK (Cached)',
      headers: {},
      config: error.config,
    });
  }
  return Promise.reject(error);
});
```

**Cache Invalidation on Mutations:**
```javascript
// After creating/updating/deleting a template
apiCache.invalidate(/\/api\/project\/\d+\/templates/);

// After task status change
apiCache.invalidate(/\/api\/project\/\d+\/tasks/);
```

**Quick Win:** ✅ Can be implemented in < 8 hours

---

### P2.3: Skeleton Loaders for All Async Content

**Current State:** Some components use `v-skeleton-loader`, but not consistently.

**Implementation:**

```vue
<!-- web/src/components/ContentSkeleton.vue -->
<template>
  <div>
    <v-skeleton-loader
      v-if="loading"
      :type="type"
      :tile="tile"
      :height="height"
    />
    <slot v-else />
  </div>
</template>

<script>
export default {
  name: 'ContentSkeleton',
  props: {
    loading: Boolean,
    type: {
      type: String,
      default: 'card',
    },
    tile: Boolean,
    height: [String, Number],
  },
};
</script>
```

**Usage:**
```vue
<ContentSkeleton :loading="tasks === null" type="table">
  <TaskList :tasks="tasks" />
</ContentSkeleton>
```

**Quick Win:** ✅ Can be implemented in < 4 hours

---

### P2.4: Optimistic UI Updates

**Current State:** UI waits for server response before updating.

**Implementation:**

```javascript
// web/src/mixins/OptimisticUpdate.js
export default {
  methods: {
    async optimisticUpdate(updateFn, apiCall, rollbackFn) {
      // Apply optimistic update
      const previousState = updateFn();
      
      try {
        // Make API call
        await apiCall();
      } catch (error) {
        // Rollback on error
        if (rollbackFn) {
          rollbackFn(previousState);
        } else {
          // Default rollback
          this.$forceUpdate();
        }
        
        // Show error
        toast.apiError(error);
        throw error;
      }
    },
  },
};
```

**Usage:**
```vue
<script>
import OptimisticUpdate from '@/mixins/OptimisticUpdate';

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
          // Rollback
          previousState.task.status = previousState.previousStatus;
        }
      );
    },
  },
};
</script>
```

**Quick Win:** ✅ Can be implemented in < 6 hours

---

## 🔮 P3: Nice-to-Have Improvements

### P3.1: Storybook for Component Documentation

**Setup:**
```bash
npx sb init --type vue
```

**Example Story:**
```javascript
// web/src/components/TaskStatus.stories.js
import TaskStatus from './TaskStatus.vue';

export default {
  title: 'Components/TaskStatus',
  component: TaskStatus,
};

export const Success = () => ({
  components: { TaskStatus },
  template: '<TaskStatus status="success" />',
});

export const Error = () => ({
  components: { TaskStatus },
  template: '<TaskStatus status="error" />',
});

export const Running = () => ({
  components: { TaskStatus },
  template: '<TaskStatus status="running" />',
});
```

**Quick Win:** ✅ Setup in < 2 hours, stories can be added incrementally

---

### P3.2: E2E Tests with Playwright

**Current State:** Basic Playwright setup exists in `e2e/` directory.

**Enhancement:** Add comprehensive E2E tests for critical flows.

**Example Test:**
```typescript
// e2e/tests/task-execution.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Task Execution Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Login
    await page.goto('/auth/login');
    await page.fill('input[name="auth"]', 'admin');
    await page.fill('input[name="password"]', 'password');
    await page.click('button[type="submit"]');
    await page.waitForURL('/project/*');
  });
  
  test('should create and run a task', async ({ page }) => {
    // Navigate to template
    await page.goto('/project/1/templates/1');
    
    // Click run button
    await page.click('button:has-text("Run")');
    
    // Fill task form
    await page.fill('input[name="message"]', 'Test task');
    await page.click('button:has-text("Start")');
    
    // Wait for task to appear in list
    await expect(page.locator('.task-list')).toContainText('Test task');
    
    // Open task logs
    await page.click('.task-list .task-row:first-child');
    await expect(page.locator('.task-log-view')).toBeVisible();
    
    // Verify log output appears
    await expect(page.locator('.task-log-records')).toContainText(/TASK|PLAY/);
  });
  
  test('should stop a running task', async ({ page }) => {
    // ... test implementation
  });
});
```

**Quick Win:** ⚠️ Moderate effort - requires test infrastructure setup

---

### P3.3: Performance Monitoring

**Implementation:**

```javascript
// web/src/lib/performance.js
export function measurePerformance(name, fn) {
  if (process.env.NODE_ENV === 'production' && 'performance' in window) {
    performance.mark(`${name}-start`);
    
    const result = fn();
    
    if (result instanceof Promise) {
      return result.finally(() => {
        performance.mark(`${name}-end`);
        performance.measure(name, `${name}-start`, `${name}-end`);
        
        const measure = performance.getEntriesByName(name)[0];
        console.log(`[Performance] ${name}: ${measure.duration.toFixed(2)}ms`);
        
        // Send to analytics
        if (window.gtag) {
          window.gtag('event', 'timing_complete', {
            name,
            value: Math.round(measure.duration),
          });
        }
      });
    } else {
      performance.mark(`${name}-end`);
      performance.measure(name, `${name}-start`, `${name}-end`);
      return result;
    }
  }
  
  return fn();
}
```

**Usage:**
```javascript
const tasks = await measurePerformance('loadTasks', () =>
  api.get('/tasks')
);
```

**Quick Win:** ✅ Can be implemented in < 4 hours

---

## 📊 Quick Wins Summary

| Task | Priority | Effort | Impact |
|------|----------|--------|--------|
| XSS Sanitization | P0 | 4h | 🔴 Critical |
| API Client with Interceptors | P0 | 6h | 🔴 Critical |
| Breadcrumb Navigation | P1 | 4h | 🟡 High |
| Confirmation Dialogs | P1 | 4h | 🟡 High |
| TypeScript Setup | P1 | 4h | 🟡 High |
| Enhanced Ansible Output Parser | P1 | 8h | 🟡 High |
| Ansible Facts Viewer | P1 | 6h | 🟡 High |
| Ansible Tags Manager | P1 | 6h | 🟡 High |
| API Caching | P2 | 8h | 🟢 Medium |
| Skeleton Loaders | P2 | 4h | 🟢 Medium |
| Inventory Visualizer | P2 | 8h | 🟢 Medium |
| Playbook Syntax Highlighter | P3 | 4h | 🟢 Low |
| Storybook Setup | P3 | 2h | 🟢 Low |

**Total Quick Wins:** ~68 hours of focused work

---

## 🎯 Implementation Roadmap

### Week 1: Critical Security (P0)
- Day 1-2: XSS sanitization + Secret masking
- Day 3-4: API client with interceptors
- Day 5: Testing & documentation

### Week 2: High-Priority UX (P1)
- Day 1-2: Breadcrumbs + Confirmations
- Day 3-4: Enhanced log viewer
- Day 5: Virtualized lists

### Week 3: Ansible Integration (P1)
- Day 1-2: Enhanced Ansible Output Parser
- Day 3: Ansible Facts Viewer
- Day 4: Ansible Tags Manager
- Day 5: Playbook Structure Visualizer (start)

### Week 4: Ansible Integration & Architecture (P1-P2)
- Day 1: Playbook Structure Visualizer (finish)
- Day 2-3: TypeScript setup + migration start
- Day 4: API caching
- Day 5: Inventory Visualizer

### Week 5: Polish & Testing (P2-P3)
- Day 1-2: Accessibility improvements
- Day 3: Skeleton loaders + Optimistic UI
- Day 4: E2E tests
- Day 5: Performance monitoring + Storybook

### Week 6-7: Advanced Ansible Features (P2-P3)
- Week 6: Ansible Lint Integration (requires backend)
- Week 7: Ansible Galaxy Integration (requires backend) + Playbook Syntax Highlighter

---

## ⚠️ Risk Factors & Mitigation

### Breaking Changes
- **Risk:** TypeScript migration may break existing code
- **Mitigation:** Gradual migration, `// @ts-ignore` for legacy code

### Performance Regression
- **Risk:** Adding sanitization may slow log rendering
- **Mitigation:** Benchmark before/after, use Web Workers for heavy sanitization

### User Experience Disruption
- **Risk:** New confirmation dialogs may annoy users
- **Mitigation:** Make confirmations optional via user settings, skip for non-destructive actions

### Browser Compatibility
- **Risk:** New features may not work in older browsers
- **Mitigation:** Use Babel transforms, polyfills, feature detection

---

---

## 🎯 P1-P2: Ansible-Specific Enhancements

### Current Ansible Integration State

**What's Already Implemented:**
- ✅ Basic Ansible output parsing (stages, hosts, errors)
- ✅ AnsibleStageView component (host statistics, failed tasks)
- ✅ TaskParamsAnsibleForm (tags, skip_tags, limit, debug, dry_run, diff)
- ✅ ANSI color support in logs
- ✅ Basic task status tracking

**What's Missing:**
- ❌ Structured Ansible facts viewer
- ❌ Playbook structure visualization
- ❌ Task dependencies graph
- ❌ Handlers display
- ❌ Variables (vars) viewer
- ❌ Ansible Galaxy integration
- ❌ Ansible Collections browser
- ❌ ansible-lint integration
- ❌ Inventory visualization
- ❌ Playbook syntax highlighting
- ❌ Ansible Vault status indicator
- ❌ Task results structured view
- ❌ Playbook execution flow diagram

---

### P1.6: Enhanced Ansible Output Parser

**Current State:** Logs are parsed only for basic task boundaries. No structured parsing of Ansible-specific output.

**Improvement:** Parse Ansible output into structured data (plays, tasks, handlers, facts, variables).

**Implementation:**

```javascript
// web/src/lib/ansibleParser.js
/**
 * Parse Ansible playbook output into structured format
 */
export class AnsibleOutputParser {
  constructor() {
    this.plays = [];
    this.currentPlay = null;
    this.currentTask = null;
    this.facts = {};
    this.handlers = [];
    this.variables = {};
  }
  
  parse(logs) {
    this.reset();
    
    logs.forEach((log) => {
      const output = log.output || '';
      
      // Parse PLAY
      const playMatch = output.match(/^PLAY\s+\[([^\]]+)\]/);
      if (playMatch) {
        this.startPlay(playMatch[1], log.time);
        return;
      }
      
      // Parse TASK
      const taskMatch = output.match(/^TASK\s+\[([^\]]+)\]/);
      if (taskMatch) {
        this.startTask(taskMatch[1], log.time);
        return;
      }
      
      // Parse task result
      const resultMatch = output.match(/^(ok|changed|failed|skipped|unreachable):\s*\[([^\]]+)\]\s*(.*)/);
      if (resultMatch) {
        this.addTaskResult(resultMatch[1], resultMatch[2], resultMatch[3], log.time);
        return;
      }
      
      // Parse facts gathering
      if (output.includes('GATHERING FACTS')) {
        this.currentTask = {
          name: 'Gathering Facts',
          type: 'setup',
          status: 'running',
          startTime: log.time,
        };
        return;
      }
      
      // Parse handlers
      const handlerMatch = output.match(/^RUNNING HANDLER\s+\[([^\]]+)\]/);
      if (handlerMatch) {
        this.startHandler(handlerMatch[1], log.time);
        return;
      }
      
      // Parse variables (vars_prompt, set_fact, etc.)
      const varMatch = output.match(/^ok:\s*\[([^\]]+)\]\s*=>\s*\{[^}]*"ansible_facts"[^}]*\}/);
      if (varMatch) {
        this.parseFacts(output);
        return;
      }
      
      // Add output to current task
      if (this.currentTask) {
        if (!this.currentTask.output) {
          this.currentTask.output = [];
        }
        this.currentTask.output.push(log);
      }
    });
    
    return {
      plays: this.plays,
      facts: this.facts,
      handlers: this.handlers,
      variables: this.variables,
    };
  }
  
  startPlay(name, time) {
    if (this.currentPlay) {
      this.plays.push(this.currentPlay);
    }
    
    this.currentPlay = {
      name,
      startTime: time,
      tasks: [],
      hosts: new Set(),
      status: 'running',
    };
    this.currentTask = null;
  }
  
  startTask(name, time) {
    if (this.currentTask) {
      this.currentPlay.tasks.push(this.currentTask);
    }
    
    this.currentTask = {
      name,
      startTime: time,
      status: 'running',
      results: [],
      output: [],
    };
  }
  
  addTaskResult(status, host, message, time) {
    if (!this.currentTask) return;
    
    this.currentPlay.hosts.add(host);
    
    this.currentTask.results.push({
      host,
      status,
      message,
      time,
    });
    
    // Update task status
    if (status === 'failed' || status === 'unreachable') {
      this.currentTask.status = 'failed';
    } else if (status === 'changed' && this.currentTask.status !== 'failed') {
      this.currentTask.status = 'changed';
    } else if (this.currentTask.status === 'running') {
      this.currentTask.status = 'ok';
    }
  }
  
  startHandler(name, time) {
    this.handlers.push({
      name,
      startTime: time,
      status: 'running',
      results: [],
    });
  }
  
  parseFacts(output) {
    try {
      // Extract JSON from output
      const jsonMatch = output.match(/\{.*"ansible_facts".*\}/);
      if (jsonMatch) {
        const facts = JSON.parse(jsonMatch[0]);
        Object.assign(this.facts, facts.ansible_facts || {});
      }
    } catch (e) {
      console.warn('Failed to parse facts:', e);
    }
  }
  
  reset() {
    this.plays = [];
    this.currentPlay = null;
    this.currentTask = null;
    this.facts = {};
    this.handlers = [];
    this.variables = {};
  }
}

export const ansibleParser = new AnsibleOutputParser();
```

**Quick Win:** ✅ Can be implemented in < 8 hours

---

### P1.7: Ansible Facts Viewer

**Current State:** Facts are not displayed in a structured way.

**Improvement:** Create a dedicated component to view Ansible facts with search and filtering.

**Implementation:**

```vue
<!-- web/src/components/AnsibleFactsViewer.vue -->
<template>
  <div class="ansible-facts-viewer">
    <v-toolbar flat dense>
      <v-text-field
        v-model="searchQuery"
        prepend-inner-icon="mdi-magnify"
        label="Search facts"
        single-line
        hide-details
        clearable
        class="mr-4"
      />
      <v-select
        v-model="factCategory"
        :items="factCategories"
        label="Category"
        clearable
        dense
        hide-details
        style="max-width: 200px;"
      />
    </v-toolbar>
    
    <v-treeview
      :items="filteredFacts"
      :search="searchQuery"
      :filter="filterFacts"
      item-key="key"
      item-text="label"
      item-children="children"
      activatable
      open-on-click
    >
      <template v-slot:prepend="{ item }">
        <v-icon v-if="item.type === 'object'">mdi-folder</v-icon>
        <v-icon v-else-if="item.type === 'array'">mdi-format-list-bulleted</v-icon>
        <v-icon v-else>mdi-code-tags</v-icon>
      </template>
      
      <template v-slot:label="{ item }">
        <div class="d-flex align-center">
          <span class="font-weight-medium mr-2">{{ item.label }}</span>
          <v-chip
            v-if="item.type"
            x-small
            color="grey"
            text-color="white"
          >
            {{ item.type }}
          </v-chip>
        </div>
      </template>
      
      <template v-slot:append="{ item }">
        <v-btn
          v-if="item.value !== undefined"
          icon
          x-small
          @click="copyToClipboard(item.value)"
        >
          <v-icon small>mdi-content-copy</v-icon>
        </v-btn>
      </template>
    </v-treeview>
    
    <v-dialog v-model="valueDialog" max-width="600">
      <v-card>
        <v-card-title>Fact Value</v-card-title>
        <v-card-text>
          <pre class="fact-value-display">{{ selectedValue }}</pre>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn text @click="valueDialog = false">Close</v-btn>
          <v-btn color="primary" @click="copyToClipboard(selectedValue)">
            Copy
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { ansibleParser } from '@/lib/ansibleParser';
import toast from '@/lib/toast';

export default {
  name: 'AnsibleFactsViewer',
  props: {
    facts: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      searchQuery: '',
      factCategory: null,
      valueDialog: false,
      selectedValue: '',
    };
  },
  computed: {
    factCategories() {
      const categories = new Set();
      Object.keys(this.facts).forEach((key) => {
        const category = key.split('_')[0];
        categories.add(category);
      });
      return Array.from(categories).map((c) => ({ text: c, value: c }));
    },
    
    structuredFacts() {
      return this.buildTree(this.facts, 'facts');
    },
    
    filteredFacts() {
      if (!this.factCategory) {
        return this.structuredFacts;
      }
      
      return this.structuredFacts.filter((item) =>
        item.key.startsWith(this.factCategory)
      );
    },
  },
  methods: {
    buildTree(obj, prefix = '') {
      const items = [];
      
      Object.keys(obj).forEach((key) => {
        const value = obj[key];
        const fullKey = prefix ? `${prefix}.${key}` : key;
        
        if (value && typeof value === 'object' && !Array.isArray(value)) {
          items.push({
            key: fullKey,
            label: key,
            type: 'object',
            children: this.buildTree(value, fullKey),
          });
        } else if (Array.isArray(value)) {
          items.push({
            key: fullKey,
            label: `${key} [${value.length}]`,
            type: 'array',
            children: value.map((item, index) => ({
              key: `${fullKey}[${index}]`,
              label: `[${index}]`,
              type: typeof item,
              value: item,
            })),
          });
        } else {
          items.push({
            key: fullKey,
            label: key,
            type: typeof value,
            value,
          });
        }
      });
      
      return items;
    },
    
    filterFacts(item, query) {
      return item.label.toLowerCase().includes(query.toLowerCase());
    },
    
    copyToClipboard(value) {
      const text = typeof value === 'string' ? value : JSON.stringify(value, null, 2);
      navigator.clipboard.writeText(text).then(() => {
        toast.success('Copied to clipboard');
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.fact-value-display {
  background: #f5f5f5;
  padding: 16px;
  border-radius: 4px;
  max-height: 400px;
  overflow: auto;
  font-family: monospace;
  font-size: 12px;
}
</style>
```

**Quick Win:** ✅ Can be implemented in < 6 hours

---

### P1.8: Playbook Structure Visualizer

**Current State:** No visualization of playbook structure (plays, tasks, roles, handlers).

**Improvement:** Create a tree/graph view of playbook structure.

**Implementation:**

```vue
<!-- web/src/components/PlaybookStructureViewer.vue -->
<template>
  <div class="playbook-structure-viewer">
    <v-tabs v-model="viewMode">
      <v-tab value="tree">Tree View</v-tab>
      <v-tab value="graph">Graph View</v-tab>
      <v-tab value="list">List View</v-tab>
    </v-tabs>
    
    <v-tabs-items v-model="viewMode">
      <v-tab-item value="tree">
        <v-treeview
          :items="playbookTree"
          item-key="id"
          item-text="name"
          item-children="children"
          activatable
          open-on-click
        >
          <template v-slot:prepend="{ item }">
            <v-icon :color="getItemColor(item.type)">
              {{ getItemIcon(item.type) }}
            </v-icon>
          </template>
          
          <template v-slot:label="{ item }">
            <div class="d-flex align-center">
              <span class="font-weight-medium">{{ item.name }}</span>
              <v-chip
                v-if="item.status"
                x-small
                :color="getStatusColor(item.status)"
                class="ml-2"
              >
                {{ item.status }}
              </v-chip>
              <v-chip
                v-if="item.tags && item.tags.length > 0"
                x-small
                color="grey"
                class="ml-2"
              >
                {{ item.tags.length }} tags
              </v-chip>
            </div>
          </template>
        </v-treeview>
      </v-tab-item>
      
      <v-tab-item value="graph">
        <div ref="graphContainer" class="playbook-graph"></div>
      </v-tab-item>
      
      <v-tab-item value="list">
        <v-list>
          <v-list-item
            v-for="item in flattenedItems"
            :key="item.id"
            @click="selectItem(item)"
          >
            <v-list-item-icon>
              <v-icon :color="getItemColor(item.type)">
                {{ getItemIcon(item.type) }}
              </v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>{{ item.name }}</v-list-item-title>
              <v-list-item-subtitle>{{ item.path }}</v-list-item-subtitle>
            </v-list-item-content>
            <v-list-item-action>
              <v-chip
                v-if="item.status"
                x-small
                :color="getStatusColor(item.status)"
              >
                {{ item.status }}
              </v-chip>
            </v-list-item-action>
          </v-list-item>
        </v-list>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script>
export default {
  name: 'PlaybookStructureViewer',
  props: {
    playbook: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      viewMode: 'tree',
    };
  },
  computed: {
    playbookTree() {
      if (!this.playbook.plays) return [];
      
      return this.playbook.plays.map((play, playIndex) => ({
        id: `play-${playIndex}`,
        name: play.name,
        type: 'play',
        status: play.status,
        children: [
          ...(play.tasks || []).map((task, taskIndex) => ({
            id: `task-${playIndex}-${taskIndex}`,
            name: task.name,
            type: 'task',
            status: task.status,
            tags: task.tags,
            path: `${play.name} > ${task.name}`,
          })),
          ...(play.roles || []).map((role, roleIndex) => ({
            id: `role-${playIndex}-${roleIndex}`,
            name: role.name,
            type: 'role',
            path: `${play.name} > ${role.name}`,
            children: (role.tasks || []).map((task, taskIndex) => ({
              id: `role-task-${playIndex}-${roleIndex}-${taskIndex}`,
              name: task.name,
              type: 'task',
              status: task.status,
              path: `${play.name} > ${role.name} > ${task.name}`,
            })),
          })),
        ],
      }));
    },
    
    flattenedItems() {
      const flatten = (items) => {
        const result = [];
        items.forEach((item) => {
          result.push(item);
          if (item.children) {
            result.push(...flatten(item.children));
          }
        });
        return result;
      };
      return flatten(this.playbookTree);
    },
  },
  methods: {
    getItemIcon(type) {
      const icons = {
        play: 'mdi-play-circle',
        task: 'mdi-checkbox-marked-circle',
        role: 'mdi-folder-star',
        handler: 'mdi-hand-wave',
      };
      return icons[type] || 'mdi-file-document';
    },
    
    getItemColor(type) {
      const colors = {
        play: 'primary',
        task: 'success',
        role: 'warning',
        handler: 'info',
      };
      return colors[type] || 'grey';
    },
    
    getStatusColor(status) {
      const colors = {
        ok: 'success',
        changed: 'warning',
        failed: 'error',
        skipped: 'grey',
        running: 'info',
      };
      return colors[status] || 'grey';
    },
    
    selectItem(item) {
      this.$emit('select', item);
    },
  },
};
</script>
```

**Quick Win:** ⚠️ Moderate effort (12-16 hours) - requires graph library integration

---

### P1.9: Ansible Tags & Skip Tags Manager

**Current State:** Tags can be set in TaskParamsAnsibleForm, but no visualization of available tags in playbook.

**Improvement:** Show available tags from playbook and allow easy selection.

**Implementation:**

```vue
<!-- web/src/components/AnsibleTagsManager.vue -->
<template>
  <div class="ansible-tags-manager">
    <v-expansion-panels>
      <v-expansion-panel>
        <v-expansion-panel-header>
          <div class="d-flex align-center">
            <v-icon class="mr-2">mdi-tag-multiple</v-icon>
            <span>Available Tags ({{ allTags.length }})</span>
          </div>
        </v-expansion-panel-header>
        <v-expansion-panel-content>
          <v-chip-group
            v-model="selectedTags"
            multiple
            column
          >
            <v-chip
              v-for="tag in allTags"
              :key="tag.name"
              :value="tag.name"
              :color="getTagColor(tag)"
              outlined
            >
              <v-icon left small>{{ getTagIcon(tag) }}</v-icon>
              {{ tag.name }}
              <span class="ml-1 text-caption">({{ tag.count }})</span>
            </v-chip>
          </v-chip-group>
        </v-expansion-panel-content>
      </v-expansion-panel>
      
      <v-expansion-panel>
        <v-expansion-panel-header>
          <div class="d-flex align-center">
            <v-icon class="mr-2">mdi-tag-off</v-icon>
            <span>Skip Tags</span>
          </div>
        </v-expansion-panel-header>
        <v-expansion-panel-content>
          <v-chip-group
            v-model="selectedSkipTags"
            multiple
            column
          >
            <v-chip
              v-for="tag in allTags"
              :key="tag.name"
              :value="tag.name"
              color="error"
              outlined
            >
              <v-icon left small>mdi-tag-off</v-icon>
              {{ tag.name }}
            </v-chip>
          </v-chip-group>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
    
    <v-divider class="my-4" />
    
    <div class="d-flex align-center">
      <v-text-field
        v-model="customTag"
        label="Add custom tag"
        outlined
        dense
        hide-details
        class="mr-2"
        @keyup.enter="addCustomTag"
      />
      <v-btn @click="addCustomTag" color="primary">Add</v-btn>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AnsibleTagsManager',
  props: {
    playbook: Object,
    value: {
      type: Array,
      default: () => [],
    },
    skipTags: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      customTag: '',
    };
  },
  computed: {
    allTags() {
      if (!this.playbook || !this.playbook.plays) return [];
      
      const tagMap = new Map();
      
      this.playbook.plays.forEach((play) => {
        (play.tasks || []).forEach((task) => {
          (task.tags || []).forEach((tag) => {
            if (!tagMap.has(tag)) {
              tagMap.set(tag, { name: tag, count: 0, tasks: [] });
            }
            tagMap.get(tag).count += 1;
            tagMap.get(tag).tasks.push(task);
          });
        });
      });
      
      return Array.from(tagMap.values()).sort((a, b) => b.count - a.count);
    },
    
    selectedTags: {
      get() {
        return this.value || [];
      },
      set(value) {
        this.$emit('input', value);
      },
    },
    
    selectedSkipTags: {
      get() {
        return this.skipTags || [];
      },
      set(value) {
        this.$emit('update:skipTags', value);
      },
    },
  },
  methods: {
    getTagColor(tag) {
      // Color based on tag usage frequency
      if (tag.count > 10) return 'primary';
      if (tag.count > 5) return 'success';
      return 'grey';
    },
    
    getTagIcon(tag) {
      // Special icons for common tags
      const iconMap = {
        always: 'mdi-repeat',
        never: 'mdi-close-circle',
        debug: 'mdi-bug',
        test: 'mdi-flask',
      };
      return iconMap[tag.name] || 'mdi-tag';
    },
    
    addCustomTag() {
      if (this.customTag && !this.selectedTags.includes(this.customTag)) {
        this.selectedTags = [...this.selectedTags, this.customTag];
        this.customTag = '';
      }
    },
  },
};
</script>
```

**Quick Win:** ✅ Can be implemented in < 6 hours

---

### P2.5: Ansible Inventory Visualizer

**Current State:** Inventory is managed but not visualized in a user-friendly way.

**Improvement:** Create a visual representation of inventory groups and hosts.

**Implementation:**

```vue
<!-- web/src/components/AnsibleInventoryViewer.vue -->
<template>
  <div class="ansible-inventory-viewer">
    <v-tabs v-model="viewMode">
      <v-tab value="tree">Tree View</v-tab>
      <v-tab value="groups">Groups</v-tab>
      <v-tab value="hosts">All Hosts</v-tab>
    </v-tabs>
    
    <v-tabs-items v-model="viewMode">
      <v-tab-item value="tree">
        <v-treeview
          :items="inventoryTree"
          item-key="id"
          item-text="name"
          item-children="children"
          activatable
          open-on-click
        >
          <template v-slot:prepend="{ item }">
            <v-icon :color="item.type === 'group' ? 'primary' : 'success'">
              {{ item.type === 'group' ? 'mdi-folder' : 'mdi-server' }}
            </v-icon>
          </template>
          
          <template v-slot:label="{ item }">
            <div class="d-flex align-center">
              <span class="font-weight-medium">{{ item.name }}</span>
              <v-chip
                v-if="item.vars && Object.keys(item.vars).length > 0"
                x-small
                color="grey"
                class="ml-2"
              >
                {{ Object.keys(item.vars).length }} vars
              </v-chip>
            </div>
          </template>
        </v-treeview>
      </v-tab-item>
      
      <v-tab-item value="groups">
        <v-list>
          <v-list-item
            v-for="group in groups"
            :key="group.name"
            @click="selectGroup(group)"
          >
            <v-list-item-icon>
              <v-icon color="primary">mdi-folder</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>{{ group.name }}</v-list-item-title>
              <v-list-item-subtitle>
                {{ group.hosts.length }} host(s)
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-tab-item>
      
      <v-tab-item value="hosts">
        <v-data-table
          :headers="hostHeaders"
          :items="allHosts"
          :search="hostSearch"
        >
          <template v-slot:top>
            <v-text-field
              v-model="hostSearch"
              label="Search hosts"
              prepend-inner-icon="mdi-magnify"
              class="mx-4"
            />
          </template>
          
          <template v-slot:item.groups="{ item }">
            <v-chip
              v-for="group in item.groups"
              :key="group"
              x-small
              class="mr-1"
            >
              {{ group }}
            </v-chip>
          </template>
        </v-data-table>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script>
export default {
  name: 'AnsibleInventoryViewer',
  props: {
    inventory: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      viewMode: 'tree',
      hostSearch: '',
      hostHeaders: [
        { text: 'Host', value: 'name' },
        { text: 'Groups', value: 'groups' },
        { text: 'Variables', value: 'vars' },
      ],
    };
  },
  computed: {
    groups() {
      return Object.keys(this.inventory).map((groupName) => ({
        name: groupName,
        hosts: this.inventory[groupName].hosts || [],
        vars: this.inventory[groupName].vars || {},
      }));
    },
    
    allHosts() {
      const hosts = [];
      this.groups.forEach((group) => {
        group.hosts.forEach((host) => {
          const existingHost = hosts.find((h) => h.name === host);
          if (existingHost) {
            existingHost.groups.push(group.name);
          } else {
            hosts.push({
              name: host,
              groups: [group.name],
              vars: {},
            });
          }
        });
      });
      return hosts;
    },
    
    inventoryTree() {
      return this.groups.map((group, index) => ({
        id: `group-${index}`,
        name: group.name,
        type: 'group',
        vars: group.vars,
        children: group.hosts.map((host, hostIndex) => ({
          id: `host-${index}-${hostIndex}`,
          name: host,
          type: 'host',
        })),
      }));
    },
  },
  methods: {
    selectGroup(group) {
      this.$emit('select-group', group);
    },
  },
};
</script>
```

**Quick Win:** ✅ Can be implemented in < 8 hours

---

### P2.6: Ansible Lint Integration

**Current State:** No integration with ansible-lint for playbook validation.

**Improvement:** Show ansible-lint warnings/errors in UI before running playbook.

**Implementation:**

```vue
<!-- web/src/components/AnsibleLintResults.vue -->
<template>
  <div class="ansible-lint-results">
    <v-alert
      v-if="lintResults.length === 0"
      type="success"
      text
    >
      No linting issues found!
    </v-alert>
    
    <v-list v-else>
      <v-list-item
        v-for="(issue, index) in lintResults"
        :key="index"
        :class="getIssueClass(issue.type)"
      >
        <v-list-item-icon>
          <v-icon :color="getIssueColor(issue.type)">
            {{ getIssueIcon(issue.type) }}
          </v-icon>
        </v-list-item-icon>
        
        <v-list-item-content>
          <v-list-item-title>{{ issue.message }}</v-list-item-title>
          <v-list-item-subtitle>
            {{ issue.rule }} - Line {{ issue.line }}:{{ issue.column }}
          </v-list-item-subtitle>
        </v-list-item-content>
        
        <v-list-item-action>
          <v-btn
            icon
            small
            @click="goToLine(issue.line)"
          >
            <v-icon small>mdi-arrow-right</v-icon>
          </v-btn>
        </v-list-item-action>
      </v-list-item>
    </v-list>
  </div>
</template>

<script>
export default {
  name: 'AnsibleLintResults',
  props: {
    lintResults: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    getIssueType(type) {
      const types = {
        error: 'error',
        warning: 'warning',
        info: 'info',
      };
      return types[type] || 'info';
    },
    
    getIssueColor(type) {
      const colors = {
        error: 'error',
        warning: 'warning',
        info: 'info',
      };
      return colors[type] || 'grey';
    },
    
    getIssueIcon(type) {
      const icons = {
        error: 'mdi-alert-circle',
        warning: 'mdi-alert',
        info: 'mdi-information',
      };
      return icons[type] || 'mdi-information';
    },
    
    getIssueClass(type) {
      return `ansible-lint-issue--${type}`;
    },
    
    goToLine(line) {
      this.$emit('go-to-line', line);
    },
  },
};
</script>

<style lang="scss" scoped>
.ansible-lint-issue--error {
  border-left: 4px solid #f44336;
}

.ansible-lint-issue--warning {
  border-left: 4px solid #ff9800;
}

.ansible-lint-issue--info {
  border-left: 4px solid #2196f3;
}
</style>
```

**Backend Integration Required:**
```javascript
// API endpoint to lint playbook
// POST /api/project/:projectId/templates/:templateId/lint
// Returns: { issues: [...] }
```

**Quick Win:** ⚠️ Requires backend support (16-24 hours total)

---

### P2.7: Ansible Galaxy Integration

**Current State:** No integration with Ansible Galaxy.

**Improvement:** Browse and install roles/collections from Ansible Galaxy.

**Implementation:**

```vue
<!-- web/src/components/AnsibleGalaxyBrowser.vue -->
<template>
  <div class="ansible-galaxy-browser">
    <v-tabs v-model="contentType">
      <v-tab value="roles">Roles</v-tab>
      <v-tab value="collections">Collections</v-tab>
    </v-tabs>
    
    <v-tabs-items v-model="contentType">
      <v-tab-item value="roles">
        <v-text-field
          v-model="searchQuery"
          prepend-inner-icon="mdi-magnify"
          label="Search roles"
          @keyup.enter="searchGalaxy"
          class="mb-4"
        />
        
        <v-list v-if="searchResults.length > 0">
          <v-list-item
            v-for="role in searchResults"
            :key="role.id"
          >
            <v-list-item-avatar>
              <v-icon>mdi-package-variant</v-icon>
            </v-list-item-avatar>
            
            <v-list-item-content>
              <v-list-item-title>{{ role.name }}</v-list-item-title>
              <v-list-item-subtitle>
                {{ role.description }}
              </v-list-item-subtitle>
              <div class="mt-2">
                <v-chip x-small class="mr-1">
                  <v-icon x-small left>mdi-star</v-icon>
                  {{ role.star_count }}
                </v-chip>
                <v-chip x-small class="mr-1">
                  <v-icon x-small left>mdi-download</v-icon>
                  {{ role.download_count }}
                </v-chip>
                <v-chip x-small>
                  <v-icon x-small left>mdi-tag</v-icon>
                  {{ role.version }}
                </v-chip>
              </div>
            </v-list-item-content>
            
            <v-list-item-action>
              <v-btn
                color="primary"
                @click="installRole(role)"
              >
                Install
              </v-btn>
            </v-list-item-action>
          </v-list-item>
        </v-list>
      </v-tab-item>
      
      <v-tab-item value="collections">
        <!-- Similar structure for collections -->
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script>
import { api } from '@/lib/apiClient';

export default {
  name: 'AnsibleGalaxyBrowser',
  props: {
    projectId: Number,
  },
  data() {
    return {
      contentType: 'roles',
      searchQuery: '',
      searchResults: [],
      loading: false,
    };
  },
  methods: {
    async searchGalaxy() {
      if (!this.searchQuery) return;
      
      this.loading = true;
      try {
        // This would require backend API integration with Ansible Galaxy API
        const response = await api.get('/ansible-galaxy/search', {
          params: {
            type: this.contentType,
            query: this.searchQuery,
          },
        });
        this.searchResults = response.data.results;
      } catch (error) {
        console.error('Failed to search Galaxy:', error);
      } finally {
        this.loading = false;
      }
    },
    
    async installRole(role) {
      try {
        await api.post(`/api/project/${this.projectId}/ansible-galaxy/install`, {
          type: 'role',
          name: role.name,
          version: role.version,
        });
        this.$emit('role-installed', role);
      } catch (error) {
        console.error('Failed to install role:', error);
      }
    },
  },
};
</script>
```

**Quick Win:** ⚠️ Requires backend API integration (24-32 hours total)

---

### P3.4: Playbook Syntax Highlighter

**Current State:** Playbook YAML is displayed as plain text.

**Improvement:** Add syntax highlighting for Ansible playbook YAML.

**Implementation:**

```vue
<!-- web/src/components/PlaybookEditor.vue -->
<template>
  <div class="playbook-editor">
    <codemirror
      v-model="playbookContent"
      :options="editorOptions"
      @ready="onEditorReady"
    />
  </div>
</template>

<script>
import { codemirror } from 'vue-codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/yaml/yaml';
import 'codemirror/addon/lint/lint';
import 'codemirror/addon/lint/yaml-lint';
import 'codemirror/addon/edit/closebrackets';
import 'codemirror/addon/edit/matchbrackets';

export default {
  name: 'PlaybookEditor',
  components: {
    codemirror,
  },
  props: {
    value: String,
  },
  data() {
    return {
      editorOptions: {
        mode: 'yaml',
        theme: 'default',
        lineNumbers: true,
        lineWrapping: true,
        indentUnit: 2,
        tabSize: 2,
        autoCloseBrackets: true,
        matchBrackets: true,
        lint: true,
        gutters: ['CodeMirror-lint-markers'],
      },
    };
  },
  computed: {
    playbookContent: {
      get() {
        return this.value;
      },
      set(value) {
        this.$emit('input', value);
      },
    },
  },
  methods: {
    onEditorReady(editor) {
      // Add Ansible-specific hints
      // This would require custom CodeMirror mode or addon
    },
  },
};
</script>
```

**Quick Win:** ✅ Can be implemented in < 4 hours (basic highlighting)

---

## 📊 Ansible Enhancements Summary

| Task | Priority | Effort | Impact |
|------|----------|--------|--------|
| Enhanced Ansible Output Parser | P1 | 8h | 🟡 High |
| Ansible Facts Viewer | P1 | 6h | 🟡 High |
| Playbook Structure Visualizer | P1 | 12-16h | 🟡 High |
| Ansible Tags Manager | P1 | 6h | 🟡 High |
| Inventory Visualizer | P2 | 8h | 🟢 Medium |
| Ansible Lint Integration | P2 | 16-24h | 🟢 Medium |
| Ansible Galaxy Integration | P2 | 24-32h | 🟢 Medium |
| Playbook Syntax Highlighter | P3 | 4h | 🟢 Low |

**Total Ansible Enhancements:** ~84-110 hours

---

## 📚 Additional Resources

- [Vue 2 → Vue 3 Migration Guide](https://v3-migration.vuejs.org/)
- [WCAG 2.1 Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
- [OWASP XSS Prevention](https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_Prevention_Cheat_Sheet.html)
- [TypeScript Vue Plugin](https://github.com/johnsoncodehk/volar)
- [Ansible Documentation](https://docs.ansible.com/)
- [Ansible Galaxy API](https://galaxy.ansible.com/api/)
- [ansible-lint Documentation](https://ansible-lint.readthedocs.io/)

---

**Last Updated:** 2024
**Review Status:** Ready for Implementation
**Estimated Total Effort:** ~204-270 hours (5-7 weeks for 1 developer)

