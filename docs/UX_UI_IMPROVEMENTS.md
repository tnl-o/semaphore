# UX/UI Improvements Guide

This document outlines UX/UI improvements for the Semaphore UI frontend.

## Overview

The frontend is built with Vue.js 2.6.14 and Vuetify 2.6.10. This guide covers improvements to enhance user experience and interface usability.

## Current State

### Existing Features

- ✅ Basic loading states (App.vue)
- ✅ Dark theme support (with localStorage persistence)
- ✅ Basic error handling (getErrorMessage utility)
- ✅ Vuetify Material Design components
- ✅ Responsive layout (partially)

### Areas for Improvement

- Loading states not consistently applied
- Error messages could be more user-friendly
- Toast notifications need improvement
- Mobile optimization needs work
- Dark theme could be enhanced

## 1. Loading States

### Current Implementation

Basic loading state exists in `App.vue`:

```vue
<v-app v-else-if="state === 'loading'">
  <v-progress-circular indeterminate></v-progress-circular>
</v-app>
```

### Improvement Plan

#### 1.1 Global Loading Mixin

Create a reusable loading mixin:

```javascript
// web/src/mixins/LoadingMixin.js
export default {
  data() {
    return {
      loading: false,
      loadingMessage: null,
    };
  },
  
  methods: {
    async withLoading(fn, message = null) {
      this.loading = true;
      this.loadingMessage = message;
      try {
        return await fn();
      } finally {
        this.loading = false;
        this.loadingMessage = null;
      }
    },
  },
};
```

#### 1.2 Component-Level Loading

Add loading indicators to components:

```vue
<template>
  <div>
    <v-overlay :value="loading" :opacity="0.7">
      <v-progress-circular
        indeterminate
        size="64"
        color="primary"
      ></v-progress-circular>
      <div v-if="loadingMessage" class="mt-4">
        {{ loadingMessage }}
      </div>
    </v-overlay>
    
    <!-- Component content -->
  </div>
</template>

<script>
import LoadingMixin from '@/mixins/LoadingMixin';

export default {
  mixins: [LoadingMixin],
  
  methods: {
    async loadData() {
      await this.withLoading(async () => {
        // Load data
      }, 'Loading data...');
    },
  },
};
</script>
```

#### 1.3 Button Loading States

Add loading states to buttons:

```vue
<v-btn
  :loading="saving"
  :disabled="saving"
  @click="save"
>
  Save
</v-btn>
```

#### 1.4 Skeleton Loaders

Use skeleton loaders for better perceived performance:

```vue
<template>
  <div v-if="loading">
    <v-skeleton-loader
      type="table"
      :loading="loading"
    ></v-skeleton-loader>
  </div>
  <div v-else>
    <!-- Actual content -->
  </div>
</template>
```

### Implementation Checklist

- [ ] Create `LoadingMixin.js`
- [ ] Add loading states to all async operations
- [ ] Add button loading states
- [ ] Implement skeleton loaders for lists
- [ ] Add progress indicators for long operations
- [ ] Test loading states on slow connections

## 2. Error Handling

### Current Implementation

Basic error handling exists in `lib/error.js`:

```javascript
export function getErrorMessage(err) {
  if (err.response?.data?.error) {
    return err.response.data.error;
  }
  return err.message;
}
```

### Improvement Plan

#### 2.1 Enhanced Error Messages

Improve error message extraction:

```javascript
// web/src/lib/error.js
export function getErrorMessage(err) {
  // Network errors
  if (!err.response) {
    if (err.code === 'ECONNABORTED') {
      return 'Request timeout. Please try again.';
    }
    if (err.message === 'Network Error') {
      return 'Network error. Please check your connection.';
    }
    return 'An unexpected error occurred. Please try again.';
  }
  
  // HTTP errors
  const status = err.response.status;
  const data = err.response.data;
  
  // Validation errors
  if (status === 400 && data.error) {
    return data.error;
  }
  
  // Authentication errors
  if (status === 401) {
    return 'Your session has expired. Please log in again.';
  }
  
  // Authorization errors
  if (status === 403) {
    return 'You do not have permission to perform this action.';
  }
  
  // Not found errors
  if (status === 404) {
    return 'The requested resource was not found.';
  }
  
  // Server errors
  if (status >= 500) {
    return 'Server error. Please try again later.';
  }
  
  // Default
  return data?.error || err.message || 'An error occurred';
}
```

#### 2.2 Toast Notifications

Create a centralized toast notification system:

```javascript
// web/src/plugins/toast.js
import Vue from 'vue';
import { EventBus } from '@/event-bus';

export default {
  install(Vue) {
    Vue.prototype.$toast = {
      success(message, options = {}) {
        EventBus.$emit('i-snackbar', {
          color: 'success',
          text: message,
          timeout: 5000,
          ...options,
        });
      },
      
      error(message, options = {}) {
        EventBus.$emit('i-snackbar', {
          color: 'error',
          text: message,
          timeout: 7000,
          ...options,
        });
      },
      
      info(message, options = {}) {
        EventBus.$emit('i-snackbar', {
          color: 'info',
          text: message,
          timeout: 5000,
          ...options,
        });
      },
      
      warning(message, options = {}) {
        EventBus.$emit('i-snackbar', {
          color: 'warning',
          text: message,
          timeout: 6000,
          ...options,
        });
      },
    };
  },
};
```

#### 2.3 Error Boundary Component

Create error boundary for better error handling:

```vue
<!-- web/src/components/ErrorBoundary.vue -->
<template>
  <div v-if="hasError">
    <v-alert
      type="error"
      prominent
      dismissible
      @input="dismiss"
    >
      <v-alert-title>Something went wrong</v-alert-title>
      {{ errorMessage }}
      <v-btn
        text
        small
        @click="retry"
        class="mt-2"
      >
        Retry
      </v-btn>
    </v-alert>
  </div>
  <slot v-else></slot>
</template>

<script>
export default {
  data() {
    return {
      hasError: false,
      errorMessage: null,
    };
  },
  
  errorCaptured(err, instance, info) {
    this.hasError = true;
    this.errorMessage = err.message || 'An unexpected error occurred';
    console.error('Error captured:', err, info);
    return false;
  },
  
  methods: {
    retry() {
      this.hasError = false;
      this.errorMessage = null;
      this.$forceUpdate();
    },
    
    dismiss() {
      this.hasError = false;
    },
  },
};
</script>
```

#### 2.4 Form Error Display

Improve form error display:

```vue
<template>
  <v-form>
    <v-text-field
      v-model="form.name"
      :error-messages="errors.name"
      label="Name"
    ></v-text-field>
    
    <v-alert
      v-if="formError"
      type="error"
      dismissible
      @input="formError = null"
      class="mb-4"
    >
      {{ formError }}
    </v-alert>
  </v-form>
</template>
```

### Implementation Checklist

- [ ] Enhance `getErrorMessage` function
- [ ] Create toast notification plugin
- [ ] Add error boundary component
- [ ] Improve form error display
- [ ] Add retry mechanisms for failed requests
- [ ] Add error logging
- [ ] Test error scenarios

## 3. Dark Theme

### Current Implementation

Dark theme is already implemented in `App.vue`:

```javascript
darkMode(val) {
  this.$vuetify.theme.dark = val;
  if (val) {
    localStorage.setItem('darkMode', '1');
  } else {
    localStorage.removeItem('darkMode');
  }
}
```

### Improvement Plan

#### 3.1 Theme Toggle Component

Create a dedicated theme toggle:

```vue
<!-- web/src/components/ThemeToggle.vue -->
<template>
  <v-btn
    icon
    @click="toggleTheme"
    :title="darkMode ? 'Switch to light theme' : 'Switch to dark theme'"
  >
    <v-icon>{{ darkMode ? 'mdi-weather-night' : 'mdi-weather-sunny' }}</v-icon>
  </v-btn>
</template>

<script>
export default {
  computed: {
    darkMode() {
      return this.$vuetify.theme.dark;
    },
  },
  
  methods: {
    toggleTheme() {
      this.$vuetify.theme.dark = !this.$vuetify.theme.dark;
      localStorage.setItem('darkMode', this.$vuetify.theme.dark ? '1' : '0');
    },
  },
};
</script>
```

#### 3.2 System Theme Detection

Detect system theme preference:

```javascript
// web/src/plugins/theme.js
export default {
  install(Vue) {
    // Detect system preference
    if (localStorage.getItem('darkMode') === null) {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      Vue.prototype.$vuetify.theme.dark = prefersDark;
    }
    
    // Listen for system theme changes
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      if (localStorage.getItem('darkMode') === null) {
        Vue.prototype.$vuetify.theme.dark = e.matches;
      }
    });
  },
};
```

#### 3.3 Theme Customization

Allow theme customization:

```javascript
// web/src/store/modules/theme.js
export default {
  namespaced: true,
  
  state: {
    dark: localStorage.getItem('darkMode') === '1',
    primaryColor: '#1976D2',
    accentColor: '#FF4081',
  },
  
  mutations: {
    SET_DARK(state, value) {
      state.dark = value;
      localStorage.setItem('darkMode', value ? '1' : '0');
    },
    
    SET_PRIMARY_COLOR(state, color) {
      state.primaryColor = color;
    },
  },
};
```

### Implementation Checklist

- [ ] Create theme toggle component
- [ ] Add system theme detection
- [ ] Improve theme persistence
- [ ] Add theme customization options
- [ ] Test theme switching
- [ ] Ensure all components support dark theme

## 4. Mobile Optimization

### Current State

Vuetify provides responsive components, but mobile experience could be improved.

### Improvement Plan

#### 4.1 Responsive Navigation

Improve mobile navigation:

```vue
<template>
  <v-navigation-drawer
    v-model="drawer"
    :temporary="$vuetify.breakpoint.mobile"
    :permanent="!$vuetify.breakpoint.mobile"
  >
    <!-- Navigation items -->
  </v-navigation-drawer>
  
  <v-app-bar>
    <v-app-bar-nav-icon
      v-if="$vuetify.breakpoint.mobile"
      @click="drawer = !drawer"
    ></v-app-bar-nav-icon>
  </v-app-bar>
</template>
```

#### 4.2 Touch-Friendly Components

Make components touch-friendly:

```vue
<v-btn
  :min-width="$vuetify.breakpoint.mobile ? 48 : 36"
  :min-height="$vuetify.breakpoint.mobile ? 48 : 36"
>
  Action
</v-btn>
```

#### 4.3 Mobile-Optimized Tables

Use cards on mobile instead of tables:

```vue
<template>
  <v-data-table
    v-if="!$vuetify.breakpoint.mobile"
    :items="items"
  ></v-data-table>
  
  <div v-else>
    <v-card
      v-for="item in items"
      :key="item.id"
      class="mb-2"
    >
      <v-card-text>
        <!-- Card content -->
      </v-card-text>
    </v-card>
  </div>
</template>
```

#### 4.4 Viewport Meta Tag

Ensure proper viewport configuration:

```html
<!-- public/index.html -->
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=5, user-scalable=yes">
```

#### 4.5 Mobile-Specific Styles

Add mobile-specific styles:

```scss
// web/src/assets/scss/mobile.scss
@media (max-width: 600px) {
  .mobile-full-width {
    width: 100% !important;
  }
  
  .mobile-padding {
    padding: 8px !important;
  }
  
  .mobile-text-small {
    font-size: 0.875rem !important;
  }
}
```

### Implementation Checklist

- [ ] Improve responsive navigation
- [ ] Make buttons touch-friendly
- [ ] Convert tables to cards on mobile
- [ ] Optimize forms for mobile
- [ ] Test on real mobile devices
- [ ] Add mobile-specific styles
- [ ] Optimize images for mobile
- [ ] Test touch interactions

## Implementation Priority

1. **High Priority:**
   - Loading states for async operations
   - Enhanced error messages
   - Toast notifications

2. **Medium Priority:**
   - Theme improvements
   - Mobile navigation
   - Touch-friendly components

3. **Low Priority:**
   - Theme customization
   - Advanced mobile optimizations
   - Skeleton loaders

## Resources

- [Vuetify Documentation](https://v2.vuetifyjs.com/)
- [Vue.js Best Practices](https://vuejs.org/style-guide/)
- [Mobile UX Best Practices](https://developers.google.com/web/fundamentals/design-and-ux/principles)
- [Accessibility Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)

