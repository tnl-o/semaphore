<template>
  <div class="ansible-facts-viewer">
    <v-toolbar flat dense>
      <v-text-field
        v-model="searchQuery"
        prepend-inner-icon="mdi-magnify"
        label="Search facts"
        single-line
        hide-details
        clearable
        class="mr-4"
      />
      <v-select
        v-model="factCategory"
        :items="factCategories"
        label="Category"
        clearable
        dense
        hide-details
        style="max-width: 200px;"
      />
    </v-toolbar>

    <v-treeview
      :items="filteredFacts"
      :search="searchQuery"
      :filter="filterFacts"
      item-key="key"
      item-text="label"
      item-children="children"
      activatable
      open-on-click
    >
      <template v-slot:prepend="{ item }">
        <v-icon v-if="item.type === 'object'">mdi-folder</v-icon>
        <v-icon v-else-if="item.type === 'array'">mdi-format-list-bulleted</v-icon>
        <v-icon v-else>mdi-code-tags</v-icon>
      </template>

      <template v-slot:label="{ item }">
        <div class="d-flex align-center">
          <span class="font-weight-medium mr-2">{{ item.label }}</span>
          <v-chip
            v-if="item.type"
            x-small
            color="grey"
            text-color="white"
          >
            {{ item.type }}
          </v-chip>
        </div>
      </template>

      <template v-slot:append="{ item }">
        <v-btn
          v-if="item.value !== undefined"
          icon
          x-small
          @click="copyToClipboard(item.value)"
        >
          <v-icon small>mdi-content-copy</v-icon>
        </v-btn>
      </template>
    </v-treeview>

    <v-dialog v-model="valueDialog" max-width="600">
      <v-card>
        <v-card-title>Fact Value</v-card-title>
        <v-card-text>
          <pre class="fact-value-display">{{ selectedValue }}</pre>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn text @click="valueDialog = false">Close</v-btn>
          <v-btn color="primary" @click="copyToClipboard(selectedValue)">
            Copy
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { toast } from '@/lib/toast';

export default {
  name: 'AnsibleFactsViewer',
  props: {
    facts: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      searchQuery: '',
      factCategory: null,
      valueDialog: false,
      selectedValue: '',
    };
  },
  computed: {
    factCategories() {
      const categories = new Set();
      Object.keys(this.facts).forEach((key) => {
        const category = key.split('_')[0];
        categories.add(category);
      });
      return Array.from(categories).map((c) => ({ text: c, value: c }));
    },

    structuredFacts() {
      return this.buildTree(this.facts, 'facts');
    },

    filteredFacts() {
      if (!this.factCategory) {
        return this.structuredFacts;
      }

      return this.structuredFacts.filter((item) => item.key.startsWith(this.factCategory));
    },
  },
  methods: {
    buildTree(obj, prefix = '') {
      const items = [];

      Object.keys(obj).forEach((key) => {
        const value = obj[key];
        const fullKey = prefix ? `${prefix}.${key}` : key;

        if (value && typeof value === 'object' && !Array.isArray(value)) {
          items.push({
            key: fullKey,
            label: key,
            type: 'object',
            children: this.buildTree(value, fullKey),
          });
        } else if (Array.isArray(value)) {
          items.push({
            key: fullKey,
            label: `${key} [${value.length}]`,
            type: 'array',
            children: value.map((item, index) => ({
              key: `${fullKey}[${index}]`,
              label: `[${index}]`,
              type: typeof item,
              value: item,
            })),
          });
        } else {
          items.push({
            key: fullKey,
            label: key,
            type: typeof value,
            value,
          });
        }
      });

      return items;
    },

    filterFacts(item, query) {
      return item.label.toLowerCase().includes(query.toLowerCase());
    },

    copyToClipboard(value) {
      const text = typeof value === 'string' ? value : JSON.stringify(value, null, 2);
      navigator.clipboard.writeText(text).then(() => {
        toast.success('Copied to clipboard');
      }).catch(() => {
        toast.error('Failed to copy to clipboard');
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.fact-value-display {
  background: #f5f5f5;
  padding: 16px;
  border-radius: 4px;
  max-height: 400px;
  overflow: auto;
  font-family: monospace;
  font-size: 12px;
}
</style>
