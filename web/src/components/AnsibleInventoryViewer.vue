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
