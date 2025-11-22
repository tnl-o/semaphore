<template>
  <v-text-field
    :value="displayValue"
    @input="handleInput"
    :type="showPassword ? 'text' : 'password'"
    :append-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
    @click:append="showPassword = !showPassword"
    v-bind="$attrs"
    v-on="$listeners"
  />
</template>

<script>
export default {
  name: 'SecurePasswordInput',
  inheritAttrs: false,
  props: {
    value: String,
  },
  data() {
    return {
      showPassword: false,
      // Never store actual value in component data
      internalValue: '',
    };
  },
  computed: {
    displayValue() {
      // Show masked value when not focused
      if (!this.showPassword && this.value) {
        return '•'.repeat(Math.min(this.value.length, 20));
      }
      return this.value || '';
    },
  },
  methods: {
    handleInput(value) {
      // Immediately emit, don't store
      this.$emit('input', value);
      // Clear internal reference after next tick
      this.$nextTick(() => {
        this.internalValue = '';
      });
    },
  },
  beforeDestroy() {
    // Clear any references
    this.internalValue = '';
  },
};
</script>
