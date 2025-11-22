<template>
  <div class="task-row">
    <v-list-item>
      <v-list-item-icon>
        <TaskStatus :status="source.status" />
      </v-list-item-icon>
      <v-list-item-content>
        <v-list-item-title>
          <TaskLink
            :task-id="source.id"
            :label="'#' + source.id"
          />
        </v-list-item-title>
        <v-list-item-subtitle>
          <span v-if="source.message">{{ source.message }}</span>
          <span v-else-if="source.commit_message">{{ source.commit_message }}</span>
        </v-list-item-subtitle>
      </v-list-item-content>
      <v-list-item-action>
        <span class="text-caption">{{ source.start | formatDate }}</span>
      </v-list-item-action>
    </v-list-item>
  </div>
</template>

<script>
import TaskStatus from './TaskStatus.vue';
import TaskLink from './TaskLink.vue';

export default {
  name: 'TaskRow',
  components: {
    TaskStatus,
    TaskLink,
  },
  props: {
    index: Number,
    source: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style scoped>
.task-row {
  border-bottom: 1px solid rgba(0, 0, 0, 0.12);
}
</style>
