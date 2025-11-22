# Vue 3 Migration Plan

This document outlines the plan for migrating Semaphore UI frontend from Vue 2 to Vue 3.

## Overview

**Current Version:** Vue 2.6.14  
**Target Version:** Vue 3.x  
**Estimated Effort:** 60 hours

## Why Migrate to Vue 3?

- **Performance**: Better performance with Composition API and improved reactivity
- **Bundle Size**: Smaller bundle size
- **TypeScript**: Better TypeScript support
- **Ecosystem**: Access to latest Vue ecosystem packages
- **Long-term Support**: Vue 2 reaches EOL in December 2023

## Breaking Changes

### 1. Global API Changes

**Vue 2:**
```javascript
import Vue from 'vue';
Vue.config.productionTip = false;
Vue.use(plugin);
```

**Vue 3:**
```javascript
import { createApp } from 'vue';
const app = createApp(App);
app.config.productionTip = false;
app.use(plugin);
```

### 2. Event Bus

**Vue 2:**
```javascript
import Vue from 'vue';
export default new Vue();
```

**Vue 3:**
```javascript
// Use mitt or custom event emitter
import mitt from 'mitt';
export default mitt();
```

### 3. Filters Removed

**Vue 2:**
```vue
{{ message | capitalize }}
```

**Vue 3:**
```vue
{{ capitalize(message) }}
```

### 4. v-model Changes

**Vue 2:**
```vue
<Component v-model="value" />
```

**Vue 3:**
```vue
<Component v-model:modelValue="value" />
```

### 5. $listeners Removed

**Vue 2:**
```vue
<Component v-on="$listeners" />
```

**Vue 3:**
```vue
<Component v-bind="$attrs" />
```

## Migration Strategy

### Phase 1: Compatibility Analysis (4 hours)

#### 1.1 Dependency Compatibility Check

**Check compatibility of:**
- [ ] Vuetify 2 → Vuetify 3
- [ ] Vue Router 3 → Vue Router 4
- [ ] Vuex 3 → Vuex 4 or Pinia
- [ ] vue-i18n 8 → vue-i18n 9
- [ ] vue-chartjs
- [ ] vue-codemirror
- [ ] vuedraggable
- [ ] vue-virtual-scroll-list

**Tools:**
- [Vue 3 Migration Guide](https://v3-migration.vuejs.org/)
- [Vuetify Migration Guide](https://v3.vuetifyjs.com/getting-started/migration-guide/)
- [Vue Router Migration Guide](https://router.vuejs.org/guide/migration/)

#### 1.2 Code Analysis

**Identify incompatible patterns:**
- [ ] Event Bus usage
- [ ] Filter usage
- [ ] $listeners usage
- [ ] Global API usage
- [ ] Vue.extend usage
- [ ] Functional components

**Script to find patterns:**
```bash
# Find EventBus usage
grep -r "EventBus" web/src

# Find filter usage
grep -r "|" web/src --include="*.vue"

# Find $listeners
grep -r "\$listeners" web/src
```

#### 1.3 Create Migration Checklist

Document all components and files that need changes.

### Phase 2: Preparation (8 hours)

#### 2.1 Update Build Tools

- [ ] Update Vue CLI or migrate to Vite
- [ ] Update build configuration
- [ ] Update TypeScript configuration (if applicable)

#### 2.2 Create Migration Branch

```bash
git checkout -b vue3-migration
```

#### 2.3 Set Up Vue 3 Environment

- [ ] Install Vue 3 dependencies
- [ ] Create test environment
- [ ] Set up compatibility build (if needed)

### Phase 3: Core Migration (20 hours)

#### 3.1 Update Entry Point

**File: `web/src/main.js`**

**Vue 2:**
```javascript
import Vue from 'vue';
import App from './App.vue';

new Vue({
  router,
  vuetify,
  i18n,
  render: h => h(App),
}).$mount('#app');
```

**Vue 3:**
```javascript
import { createApp } from 'vue';
import App from './App.vue';

const app = createApp(App);
app.use(router);
app.use(vuetify);
app.use(i18n);
app.mount('#app');
```

#### 3.2 Replace Event Bus

**File: `web/src/event-bus.js`**

**Option 1: Use mitt**
```javascript
import mitt from 'mitt';
export default mitt();
```

**Option 2: Use provide/inject**
```javascript
// In App.vue
provide('eventBus', eventBus);

// In components
inject('eventBus')
```

#### 3.3 Update Router

**File: `web/src/router/index.js`**

**Vue 2:**
```javascript
import VueRouter from 'vue-router';
Vue.use(VueRouter);
```

**Vue 3:**
```javascript
import { createRouter, createWebHistory } from 'vue-router';
export default createRouter({
  history: createWebHistory(),
  routes: [...]
});
```

#### 3.4 Update Vuetify

**File: `web/src/plugins/vuetify.js`**

**Vue 2:**
```javascript
import Vue from 'vue';
import Vuetify from 'vuetify/lib';
Vue.use(Vuetify);
```

**Vue 3:**
```javascript
import { createVuetify } from 'vuetify';
export default createVuetify({
  // Vuetify 3 configuration
});
```

#### 3.5 Update Vuex (or migrate to Pinia)

**Option 1: Vuex 4**
```javascript
import { createStore } from 'vuex';
export default createStore({
  // Store configuration
});
```

**Option 2: Pinia (Recommended)**
```javascript
import { createPinia } from 'pinia';
export default createPinia();
```

### Phase 4: Component Migration (20 hours)

#### 4.1 Remove Filters

Replace all filter usage with computed properties or methods:

**Before:**
```vue
{{ message | capitalize }}
```

**After:**
```vue
{{ capitalize(message) }}
```

#### 4.2 Update v-model

Update components using v-model:

**Before:**
```vue
<Component v-model="value" />
```

**After:**
```vue
<Component v-model:modelValue="value" />
```

#### 4.3 Replace $listeners

**Before:**
```vue
<Component v-on="$listeners" />
```

**After:**
```vue
<Component v-bind="$attrs" />
```

#### 4.4 Update Component Syntax

Update components to use Vue 3 syntax:

- [ ] Update Options API components
- [ ] Consider Composition API for new components
- [ ] Update functional components
- [ ] Update async components

### Phase 5: Testing (16 hours)

#### 5.1 Unit Tests

- [ ] Update test setup
- [ ] Update test utilities
- [ ] Fix broken tests
- [ ] Add new tests

#### 5.2 Integration Tests

- [ ] Update E2E tests
- [ ] Test all user flows
- [ ] Test on different browsers

#### 5.3 Manual Testing

- [ ] Test all pages
- [ ] Test all forms
- [ ] Test all dialogs
- [ ] Test navigation
- [ ] Test authentication
- [ ] Test task execution

### Phase 6: Optimization (8 hours)

#### 6.1 Performance

- [ ] Optimize bundle size
- [ ] Use Composition API where beneficial
- [ ] Optimize re-renders

#### 6.2 Code Quality

- [ ] Remove deprecated code
- [ ] Update documentation
- [ ] Code review

## Migration Checklist

### Dependencies
- [ ] Vue 3.x
- [ ] Vuetify 3.x
- [ ] Vue Router 4.x
- [ ] Vuex 4.x or Pinia
- [ ] vue-i18n 9.x
- [ ] Update all Vue ecosystem packages

### Core Files
- [ ] `main.js` - Update to createApp
- [ ] `router/index.js` - Update to createRouter
- [ ] `plugins/vuetify.js` - Update to createVuetify
- [ ] `event-bus.js` - Replace with mitt or provide/inject
- [ ] `store/index.js` - Update to createStore or Pinia

### Components
- [ ] Remove all filters
- [ ] Update all v-model usage
- [ ] Replace $listeners with $attrs
- [ ] Update global API usage
- [ ] Update component registration

### Testing
- [ ] Update test setup
- [ ] Fix unit tests
- [ ] Fix integration tests
- [ ] Manual testing

## Risks and Mitigation

### Risk 1: Breaking Changes in Dependencies

**Mitigation:**
- Test each dependency upgrade separately
- Use compatibility builds where available
- Have rollback plan

### Risk 2: Large Codebase

**Mitigation:**
- Migrate incrementally
- Use automated tools where possible
- Test thoroughly at each step

### Risk 3: Third-party Components

**Mitigation:**
- Check for Vue 3 compatible versions
- Find alternatives if needed
- Consider custom implementations

## Resources

- [Vue 3 Migration Guide](https://v3-migration.vuejs.org/)
- [Vuetify 3 Migration Guide](https://v3.vuetifyjs.com/getting-started/migration-guide/)
- [Vue Router 4 Migration](https://router.vuejs.org/guide/migration/)
- [Vuex 4 Migration](https://vuex.vuejs.org/guide/migrating-to-4-0-from-3-x.html)
- [Pinia Migration Guide](https://pinia.vuejs.org/cookbook/migration-vuex.html)

## Timeline

- **Week 1-2**: Compatibility analysis and preparation
- **Week 3-4**: Core migration
- **Week 5-6**: Component migration
- **Week 7-8**: Testing and optimization

## Notes

- Consider using Vue 3 Composition API for new features
- Pinia is recommended over Vuex 4 for new projects
- Vite is recommended over Vue CLI for better performance
- Test thoroughly before production deployment

