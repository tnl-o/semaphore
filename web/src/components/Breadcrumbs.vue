<template>
  <nav aria-label="Breadcrumb navigation">
    <v-breadcrumbs
      :items="breadcrumbItems"
      class="pa-0"
      large
      role="navigation"
      aria-label="Breadcrumb"
    >
      <template v-slot:divider>
        <v-icon>mdi-chevron-right</v-icon>
      </template>

      <template v-slot:item="{ item }">
        <v-breadcrumbs-item
          :to="item.to"
          :disabled="item.disabled"
          :exact="item.exact"
        >
          <v-icon v-if="item.icon" small class="mr-1">{{ item.icon }}</v-icon>
          {{ item.text }}
        </v-breadcrumbs-item>
      </template>
    </v-breadcrumbs>
  </nav>
</template>

<script>
export default {
  name: 'Breadcrumbs',
  computed: {
    breadcrumbItems() {
      const route = this.$route;
      const items = [];

      // Home
      items.push({
        text: this.$t('dashboard') || 'Dashboard',
        to: '/',
        icon: 'mdi-home',
        exact: true,
      });

      // Project level
      if (route.params.projectId) {
        // Try to get project name from store or use ID
        const projectId = parseInt(route.params.projectId, 10);
        const project = this.$store?.state?.projects?.find(
          (p) => p.id === projectId,
        ) || { name: `Project ${projectId}` };

        items.push({
          text: project.name,
          to: `/project/${route.params.projectId}`,
          icon: 'mdi-folder',
        });

        // Template level
        if (route.params.templateId) {
          const templateId = parseInt(route.params.templateId, 10);
          const template = this.$store?.state?.templates?.find(
            (t) => t.id === templateId,
          ) || { name: `Template ${templateId}` };

          items.push({
            text: template.name,
            to: route.path,
            icon: 'mdi-file-document',
            disabled: true,
          });
        }

        // Other project sections
        const sectionMap = {
          templates: { text: this.$t('templates') || 'Templates', icon: 'mdi-file-document-outline' },
          history: { text: this.$t('history') || 'History', icon: 'mdi-history' },
          settings: { text: this.$t('settings') || 'Settings', icon: 'mdi-cog' },
          team: { text: this.$t('team') || 'Team', icon: 'mdi-account-group' },
          keys: { text: this.$t('keys') || 'Keys', icon: 'mdi-key' },
          inventory: { text: this.$t('inventory') || 'Inventory', icon: 'mdi-database' },
          environment: { text: this.$t('environment') || 'Environment', icon: 'mdi-code-braces' },
          secret_storages: { text: this.$t('secretStorages') || 'Secret Storages', icon: 'mdi-lock' },
          integrations: { text: this.$t('integrations') || 'Integrations', icon: 'mdi-puzzle' },
          repositories: { text: this.$t('repositories') || 'Repositories', icon: 'mdi-source-branch' },
          schedule: { text: this.$t('schedule') || 'Schedule', icon: 'mdi-calendar-clock' },
          runners: { text: this.$t('runners') || 'Runners', icon: 'mdi-run' },
          stats: { text: this.$t('stats') || 'Stats', icon: 'mdi-chart-line' },
          activity: { text: this.$t('activity') || 'Activity', icon: 'mdi-pulse' },
        };

        const section = Object.keys(sectionMap).find(
          (s) => route.path.includes(`/${s}`),
        );

        if (section && !route.params.templateId) {
          items.push({
            ...sectionMap[section],
            to: route.path,
            disabled: true,
          });
        }
      }

      // Global sections
      if (route.path.startsWith('/users')) {
        items.push({
          text: this.$t('users') || 'Users',
          to: route.path,
          icon: 'mdi-account-multiple',
          disabled: true,
        });
      } else if (route.path.startsWith('/runners')) {
        items.push({
          text: this.$t('runners') || 'Runners',
          to: route.path,
          icon: 'mdi-run',
          disabled: true,
        });
      } else if (route.path.startsWith('/tasks')) {
        items.push({
          text: this.$t('tasks') || 'Tasks',
          to: route.path,
          icon: 'mdi-format-list-checks',
          disabled: true,
        });
      } else if (route.path.startsWith('/apps')) {
        items.push({
          text: this.$t('apps') || 'Apps',
          to: route.path,
          icon: 'mdi-apps',
          disabled: true,
        });
      } else if (route.path.startsWith('/tokens')) {
        items.push({
          text: this.$t('tokens') || 'Tokens',
          to: route.path,
          icon: 'mdi-key-variant',
          disabled: true,
        });
      }

      return items;
    },
  },
};
</script>
