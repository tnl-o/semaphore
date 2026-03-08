<template>
  <div>
    <v-toolbar flat>
      <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>Example Plugin</v-toolbar-title>
      <v-spacer></v-spacer>
    </v-toolbar>

    <v-divider/>

    <v-container>
      <v-row>
        <v-col cols="12">
          <v-card>
            <v-card-title>
              Example Plugin Page
            </v-card-title>
            <v-card-text>
              <p>This is an example plugin page.</p>
              <p>Plugin Name: {{ pluginName }}</p>
              <p>Plugin Version: {{ pluginVersion }}</p>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import EventBus from '@/event-bus';

export default {
  name: 'ExamplePlugin',
  data() {
    return {
      pluginName: 'example',
      pluginVersion: '1.0.0',
    };
  },
  methods: {
    showDrawer() {
      EventBus.$emit('i-show-drawer');
    },
  },
};
</script>

<style scoped>
</style>

