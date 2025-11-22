<template>
  <div class="playbook-editor">
    <codemirror
      v-model="playbookContent"
      :options="editorOptions"
      @ready="onEditorReady"
    />
  </div>
</template>

<script>
import { codemirror } from 'vue-codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/yaml/yaml';
import 'codemirror/addon/edit/closebrackets';
import 'codemirror/addon/edit/matchbrackets';

export default {
  name: 'PlaybookEditor',
  components: {
    codemirror,
  },
  props: {
    value: String,
  },
  data() {
    return {
      editorOptions: {
        mode: 'yaml',
        theme: 'default',
        lineNumbers: true,
        lineWrapping: true,
        indentUnit: 2,
        tabSize: 2,
        autoCloseBrackets: true,
        matchBrackets: true,
      },
    };
  },
  computed: {
    playbookContent: {
      get() {
        return this.value;
      },
      set(value) {
        this.$emit('input', value);
      },
    },
  },
  methods: {
    onEditorReady() {
      // Add Ansible-specific hints if needed
      // This would require custom CodeMirror mode or addon
    },
  },
};
</script>

<style lang="scss" scoped>
.playbook-editor {
  height: 100%;

  ::v-deep .CodeMirror {
    height: 100%;
    font-family: 'Courier New', monospace;
    font-size: 14px;
  }
}
</style>
