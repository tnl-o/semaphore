<template>
  <div class="ansible-log-viewer">
    <div
      v-for="task in parsedTasks"
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
        <TaskStatus :status="task.status" />
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
import TaskStatus from './TaskStatus.vue';

export default {
  name: 'AnsibleLogViewer',
  components: {
    SafeLogOutput,
    TaskStatus,
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
