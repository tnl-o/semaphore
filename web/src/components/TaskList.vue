<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <ContentSkeleton :loading="tasks === null" type="table">

    <NewTaskDialog
      v-model="newTaskDialog"
      :project-id="template.project_id"
      :template="template"
      :source-task="sourceTask"
    />

    <v-data-table
        :headers="headers"
        :items="tasks"
        :hide-default-footer="hideFooter"
        :footer-props="{ itemsPerPageOptions: [20] }"
        class="mt-0 TaskListTable"
    >
      <template v-slot:item.id="{ item }">
        <TaskLink
            :task-id="item.id"
            :label="'#' + item.id"
        />
        <div style="font-size: 14px;">
          <span v-if="item.message">
            <v-icon x-small>mdi-message-outline</v-icon> {{ item.message }}
          </span>
          <span v-else-if="item.commit_hash">
            <v-icon x-small>mdi-source-fork</v-icon> {{ item.commit_message }}
          </span>
        </div>
      </template>

      <template v-slot:item.version="{ item }">
        <div v-if="item.tpl_type !== ''">
          <TaskLink
              :disabled="item.tpl_type === 'build'"
              :task-id="item.build_task_id"
              :tooltip="item.tpl_type === 'build' ? item.message : (item.build_task || {}).message"
              :label="item.tpl_type === 'build' ? item.version : (item.build_task || {}).version"
              :status="item.status"
          />
        </div>
        <div v-else>&mdash;</div>
      </template>

      <template v-slot:item.status="{ item }">
        <TaskStatus :status="item.status"/>
      </template>

      <template v-slot:item.start="{ item }">
        {{ item.start | formatDate }}
      </template>

      <template v-slot:item.end="{ item }">
        {{ [item.start, item.end] | formatMilliseconds }}
      </template>

      <template v-slot:item.actions="{ item }">
        <v-btn-toggle dense :value-comparator="() => false">
      <v-btn @click="createTask(item)" :aria-label="`Re-run task ${item.id}`">
        <v-icon aria-hidden="true">mdi-replay</v-icon>
      </v-btn>
        </v-btn-toggle>
      </template>
    </v-data-table>
  </ContentSkeleton>
</template>
<style lang="scss">
.TaskListTable td {
  height: 60px !important;
}
</style>
<script>
import { api } from '@/lib/apiClient';
import TaskStatus from '@/components/TaskStatus.vue';
import TaskLink from '@/components/TaskLink.vue';
import { TEMPLATE_TYPE_ACTION_TITLES, TEMPLATE_TYPE_ICONS } from '@/lib/constants';
import NewTaskDialog from '@/components/NewTaskDialog.vue';
import ContentSkeleton from '@/components/ContentSkeleton.vue';

export default {
  components: {
    NewTaskDialog,
    TaskStatus,
    TaskLink,
    ContentSkeleton,
  },
  props: {
    template: Object,
    limit: Number,
    hideFooter: Boolean,
  },
  data() {
    return {
      headers: [
        {
          text: this.$i18n.t('taskId'),
          value: 'id',
          sortable: false,
        },
        {
          value: 'actions',
          sortable: false,
          width: '0%',
        },
        {
          text: this.$i18n.t('version'),
          value: 'version',
          sortable: false,
        },
        {
          text: this.$i18n.t('status'),
          value: 'status',
          sortable: false,
        },
        {
          text: this.$i18n.t('user'),
          value: 'user_name',
          sortable: false,
        },
        {
          text: this.$i18n.t('start'),
          value: 'start',
          sortable: false,
        },
        {
          text: this.$i18n.t('duration'),
          value: 'end',
          sortable: false,
        },
      ],
      tasks: null,
      taskId: null,
      newTaskDialog: null,
      sourceTask: null,
      TEMPLATE_TYPE_ICONS,
    };
  },
  watch: {
    async template() {
      await this.loadData();
    },
  },
  async created() {
    await this.loadData();
  },
  methods: {
    async loadData() {
      this.tasks = null;
      const response = await api.get(
        `/api/project/${this.template.project_id}/templates/${this.template.id}/tasks/last?limit=${this.limit || 200}`,
      );
      this.tasks = response.data;
    },

    getActionButtonTitle() {
      return this.$i18n.t(`Re${TEMPLATE_TYPE_ACTION_TITLES[this.template.type]}`);
    },

    createTask(task) {
      this.sourceTask = task;
      this.newTaskDialog = true;
    },
  },
};
</script>
