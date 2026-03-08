<template>
  <div class="playbook-structure-viewer">
    <v-tabs v-model="viewMode">
      <v-tab value="tree">Tree View</v-tab>
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
