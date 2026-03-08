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
  name: 'VirtualizedTaskList',
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
        { text: 'Waiting', value: 'waiting' },
      ],
    };
  },
  computed: {
    filteredTasks() {
      let filtered = this.tasks;

      // Search filter
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase();
        filtered = filtered.filter((task) => (
          task.id.toString().includes(query)
            || (task.message || '').toLowerCase().includes(query)
            || (task.commit_message || '').toLowerCase().includes(query)
        ));
      }

      // Status filter
      if (this.statusFilter.length > 0) {
        filtered = filtered.filter((task) => this.statusFilter.includes(task.status));
      }

      return filtered;
    },
  },
};
</script>
