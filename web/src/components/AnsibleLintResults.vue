<template>
  <div class="ansible-lint-results">
    <v-alert
      v-if="lintResults.length === 0"
      type="success"
      text
    >
      No linting issues found!
    </v-alert>

    <v-list v-else>
      <v-list-item
        v-for="(issue, index) in lintResults"
        :key="index"
        :class="getIssueClass(issue.type)"
      >
        <v-list-item-icon>
          <v-icon :color="getIssueColor(issue.type)">
            {{ getIssueIcon(issue.type) }}
          </v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>{{ issue.message }}</v-list-item-title>
          <v-list-item-subtitle>
            {{ issue.rule }} - Line {{ issue.line }}:{{ issue.column }}
          </v-list-item-subtitle>
        </v-list-item-content>

        <v-list-item-action>
          <v-btn
            icon
            small
            @click="goToLine(issue.line)"
          >
            <v-icon small>mdi-arrow-right</v-icon>
          </v-btn>
        </v-list-item-action>
      </v-list-item>
    </v-list>
  </div>
</template>

<script>
export default {
  name: 'AnsibleLintResults',
  props: {
    lintResults: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    getIssueColor(type) {
      const colors = {
        error: 'error',
        warning: 'warning',
        info: 'info',
      };
      return colors[type] || 'grey';
    },

    getIssueIcon(type) {
      const icons = {
        error: 'mdi-alert-circle',
        warning: 'mdi-alert',
        info: 'mdi-information',
      };
      return icons[type] || 'mdi-information';
    },

    getIssueClass(type) {
      return `ansible-lint-issue--${type}`;
    },

    goToLine(line) {
      this.$emit('go-to-line', line);
    },
  },
};
</script>

<style lang="scss" scoped>
.ansible-lint-issue--error {
  border-left: 4px solid #f44336;
}

.ansible-lint-issue--warning {
  border-left: 4px solid #ff9800;
}

.ansible-lint-issue--info {
  border-left: 4px solid #2196f3;
}
</style>
