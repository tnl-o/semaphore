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
