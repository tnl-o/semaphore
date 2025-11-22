# TypeScript Migration Plan

This document outlines the plan for adding TypeScript support to the Semaphore UI frontend.

## Overview

**Current State:** JavaScript (ES6+)  
**Target:** TypeScript  
**Estimated Effort:** 50 hours

## Why TypeScript?

- **Type Safety**: Catch errors at compile time
- **Better IDE Support**: Improved autocomplete and refactoring
- **Documentation**: Types serve as documentation
- **Refactoring**: Safer refactoring with type checking
- **Team Collaboration**: Better code understanding

## Migration Strategy

### Phase 1: Setup (8 hours)

#### 1.1 Install TypeScript

```bash
cd web
npm install --save-dev typescript @typescript-eslint/parser @typescript-eslint/eslint-plugin
```

#### 1.2 Create TypeScript Configuration

**File: `web/tsconfig.json`**

```json
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
    }
  },
  "include": [
    "src/**/*.ts",
    "src/**/*.tsx",
    "src/**/*.vue"
  ],
  "exclude": [
    "node_modules"
  ]
}
```

#### 1.3 Update Vue Configuration

**File: `web/vue.config.js`**

```javascript
module.exports = {
  configureWebpack: {
    resolve: {
      extensions: ['.ts', '.tsx', '.js', '.vue', '.json']
    },
    module: {
      rules: [
        {
          test: /\.tsx?$/,
          loader: 'ts-loader',
          options: {
            appendTsSuffixTo: [/\.vue$/]
          }
        }
      ]
    }
  }
};
```

#### 1.4 Add Type Definitions for Vue

**File: `web/src/shims-vue.d.ts`**

```typescript
declare module '*.vue' {
  import { DefineComponent } from 'vue';
  const component: DefineComponent<{}, {}, any>;
  export default component;
}
```

### Phase 2: Gradual Migration (30 hours)

#### 2.1 Start with Utilities

Migrate utility files first (low risk):

**File: `web/src/lib/error.ts`**

```typescript
import { AxiosError } from 'axios';

export function getErrorMessage(err: unknown): string {
  if (err && typeof err === 'object' && 'response' in err) {
    const axiosError = err as AxiosError<{ error?: string }>;
    
    if (axiosError.response?.data?.error) {
      return axiosError.response.data.error;
    }
    
    if (axiosError.message && !axiosError.message.startsWith('Request failed')) {
      return axiosError.message;
    }
    
    switch (axiosError.response?.status) {
      case 401:
        return `${axiosError.response.status} ${axiosError.response.statusText}`;
      default:
        return axiosError.message || 'An error occurred';
    }
  }
  
  if (err instanceof Error) {
    return err.message;
  }
  
  return 'An unexpected error occurred';
}
```

#### 2.2 Type API Client

**File: `web/src/lib/api.ts`**

```typescript
import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

export interface ApiResponse<T = any> {
  data: T;
  status: number;
  statusText: string;
}

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: document.baseURI,
      withCredentials: true,
    });
  }

  async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.get<T>(url, config);
    return {
      data: response.data,
      status: response.status,
      statusText: response.statusText,
    };
  }

  async post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.post<T>(url, data, config);
    return {
      data: response.data,
      status: response.status,
      statusText: response.statusText,
    };
  }

  // Add other methods...
}

export default new ApiClient();
```

#### 2.3 Type Models

**File: `web/src/types/models.ts`**

```typescript
export interface User {
  id: number;
  username: string;
  name: string;
  email: string;
  admin: boolean;
  external: boolean;
  alert: boolean;
  pro: boolean;
  created: string;
}

export interface Project {
  id: number;
  name: string;
  created: string;
  alert: boolean;
}

export interface Task {
  id: number;
  project_id: number;
  template_id: number;
  status: TaskStatus;
  created: string;
  // ... other fields
}

export enum TaskStatus {
  Waiting = 'waiting',
  Running = 'running',
  Success = 'success',
  Error = 'error',
  Stopped = 'stopped',
}
```

#### 2.4 Migrate Components Gradually

Start with simple components:

**File: `web/src/components/ThemeToggle.vue`**

```vue
<template>
  <v-btn
    icon
    @click="toggleTheme"
    :title="darkMode ? 'Switch to light theme' : 'Switch to dark theme'"
  >
    <v-icon>{{ darkMode ? 'mdi-weather-night' : 'mdi-weather-sunny' }}</v-icon>
  </v-btn>
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue';
import { useVuetify } from '@/composables/useVuetify';

export default defineComponent({
  name: 'ThemeToggle',
  setup() {
    const { darkMode, toggleTheme } = useVuetify();
    
    return {
      darkMode,
      toggleTheme,
    };
  },
});
</script>
```

### Phase 3: Advanced Typing (12 hours)

#### 3.1 Type Vuex Store

**File: `web/src/store/types.ts`**

```typescript
export interface RootState {
  user: User | null;
  projects: Project[];
  // ... other state
}

export interface UserState {
  user: User | null;
}

export interface ProjectState {
  projects: Project[];
  currentProject: Project | null;
}
```

#### 3.2 Type Router

**File: `web/src/router/types.ts`**

```typescript
import { RouteRecordRaw } from 'vue-router';

export interface RouteMeta {
  requiresAuth?: boolean;
  requiresAdmin?: boolean;
  title?: string;
}

export type AppRouteRecordRaw = RouteRecordRaw & {
  meta?: RouteMeta;
};
```

#### 3.3 Type Composables

**File: `web/src/composables/useApi.ts`**

```typescript
import { ref, Ref } from 'vue';
import api from '@/lib/api';

export function useApi<T>() {
  const data: Ref<T | null> = ref(null);
  const loading = ref(false);
  const error: Ref<Error | null> = ref(null);

  async function fetch(url: string) {
    loading.value = true;
    error.value = null;
    try {
      const response = await api.get<T>(url);
      data.value = response.data;
    } catch (err) {
      error.value = err instanceof Error ? err : new Error('Unknown error');
    } finally {
      loading.value = false;
    }
  }

  return {
    data,
    loading,
    error,
    fetch,
  };
}
```

## Migration Checklist

### Setup
- [ ] Install TypeScript and dependencies
- [ ] Create `tsconfig.json`
- [ ] Update `vue.config.js`
- [ ] Add type definitions for Vue
- [ ] Configure ESLint for TypeScript

### Utilities
- [ ] Migrate `lib/error.js` → `lib/error.ts`
- [ ] Migrate `lib/api.js` → `lib/api.ts`
- [ ] Migrate `lib/constants.js` → `lib/constants.ts`
- [ ] Add type definitions for utilities

### Types
- [ ] Create `types/models.ts`
- [ ] Create `types/api.ts`
- [ ] Create `types/store.ts`
- [ ] Create `types/router.ts`

### Components
- [ ] Migrate simple components first
- [ ] Migrate form components
- [ ] Migrate complex components
- [ ] Add prop types to all components

### Store
- [ ] Type Vuex state
- [ ] Type Vuex mutations
- [ ] Type Vuex actions
- [ ] Type Vuex getters

### Router
- [ ] Type route definitions
- [ ] Type route meta
- [ ] Type route params

### Testing
- [ ] Update test setup for TypeScript
- [ ] Add type checking to tests
- [ ] Fix type errors in tests

## Best Practices

### 1. Use Strict Mode

Enable strict mode in `tsconfig.json`:
```json
{
  "compilerOptions": {
    "strict": true
  }
}
```

### 2. Avoid `any`

Use `unknown` instead of `any` when type is truly unknown:
```typescript
function handleError(err: unknown) {
  if (err instanceof Error) {
    console.error(err.message);
  }
}
```

### 3. Use Type Assertions Sparingly

Prefer type guards:
```typescript
// Good
function isUser(obj: unknown): obj is User {
  return typeof obj === 'object' && obj !== null && 'id' in obj;
}

// Avoid
const user = obj as User;
```

### 4. Document Complex Types

```typescript
/**
 * Represents a task in the system
 * @interface Task
 */
export interface Task {
  /** Unique task identifier */
  id: number;
  /** Task status */
  status: TaskStatus;
}
```

## Resources

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [Vue 3 TypeScript Guide](https://vuejs.org/guide/typescript/overview.html)
- [TypeScript with Vue 3](https://v3.vuejs.org/guide/typescript-support.html)
- [Vuex TypeScript Guide](https://vuex.vuejs.org/guide/typescript-support.html)

## Timeline

- **Week 1**: Setup and configuration
- **Week 2-3**: Migrate utilities and types
- **Week 4-5**: Migrate components
- **Week 6**: Advanced typing and testing

## Notes

- Migrate gradually, file by file
- Keep JavaScript files working during migration
- Use `// @ts-ignore` sparingly and document why
- Add types for new code from the start
- Consider using JSDoc comments for gradual typing

