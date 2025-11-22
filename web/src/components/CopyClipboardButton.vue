<template>
  <span>
    <!-- eslint-disable-next-line vuejs-accessibility/form-control-has-label -->
    <textarea
      style="position: absolute; left: -9999px; top: -9999px;"
      ref="copy_to_clipboard_textarea"
    ></textarea>

    <v-btn
      icon
      @click="copy()"
      :large="large"
    >
      <v-icon>mdi-content-copy</v-icon>
    </v-btn>
  </span>
</template>

<script>
export default {
  props: {
    text: String,
    successMessage: {
      type: String,
      default: 'Text copied to clipboard!',
    },
    large: Boolean,
    color: String,
  },
  methods: {
    async copy() {
      try {
        const el = this.$refs.copy_to_clipboard_textarea;
        el.value = this.text;
        el.focus();
        el.select();
        const successful = document.execCommand('copy');

        if (!successful) {
          throw new Error('Fallback copy failed');
        }

        const { toast } = await import('@/lib/toast');
        toast.success(this.successMessage);
      } catch (e) {
        const { toast } = await import('@/lib/toast');
        toast.error(`Can't copy to clipboard: ${e.message}`);
      }
    },
  },
};
</script>
