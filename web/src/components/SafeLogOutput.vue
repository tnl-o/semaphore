<template>
  <div
    class="task-log-records__output"
    v-html="sanitizedOutput"
  />
</template>

<script>
import { sanitizeHtml } from '@/lib/sanitize';
import { AnsiUp } from 'ansi_up';

// Configure ANSI converter (same as in main.js)
const ansiConverter = new AnsiUp();
ansiConverter.ansi_colors = [
  [
    { rgb: [85, 85, 85], class_name: 'ansi-black' },
    { rgb: [170, 0, 0], class_name: 'ansi-red' },
    { rgb: [0, 170, 0], class_name: 'ansi-green' },
    { rgb: [255, 204, 102], class_name: 'ansi-yellow' },
    { rgb: [33, 150, 243], class_name: 'ansi-blue' },
    { rgb: [170, 0, 170], class_name: 'ansi-magenta' },
    { rgb: [0, 170, 170], class_name: 'ansi-cyan' },
    { rgb: [170, 170, 170], class_name: 'ansi-white' },
  ],
  [
    { rgb: [85, 85, 85], class_name: 'ansi-bright-black' },
    { rgb: [255, 85, 85], class_name: 'ansi-bright-red' },
    { rgb: [85, 255, 85], class_name: 'ansi-bright-green' },
    { rgb: [255, 255, 85], class_name: 'ansi-bright-yellow' },
    { rgb: [85, 85, 255], class_name: 'ansi-bright-blue' },
    { rgb: [255, 85, 255], class_name: 'ansi-bright-magenta' },
    { rgb: [85, 255, 255], class_name: 'ansi-bright-cyan' },
    { rgb: [255, 255, 255], class_name: 'ansi-bright-white' },
  ],
];

/**
 * Safe log output component
 * Renders ANSI-colored log output with XSS protection
 */
export default {
  name: 'SafeLogOutput',
  props: {
    output: {
      type: String,
      required: true,
    },
  },
  computed: {
    sanitizedOutput() {
      if (!this.output) return '';

      // Convert ANSI codes to HTML
      const html = ansiConverter.ansi_to_html(String(this.output));

      // Sanitize to prevent XSS
      return sanitizeHtml(html);
    },
  },
};
</script>

<style scoped>
.task-log-records__output {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
