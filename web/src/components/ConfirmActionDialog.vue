<template>
  <v-dialog
    :value="value"
    @input="$emit('input', $event)"
    max-width="500"
    persistent
    role="alertdialog"
    :aria-labelledby="`dialog-title-${_uid}`"
    :aria-describedby="`dialog-description-${_uid}`"
  >
    <v-card>
      <v-card-title class="headline" :id="`dialog-title-${_uid}`">
        <v-icon :color="type" class="mr-2" aria-hidden="true">{{ icon }}</v-icon>
        {{ title }}
      </v-card-title>

      <v-card-text :id="`dialog-description-${_uid}`">
        <p>{{ message }}</p>

        <!-- Additional warning for destructive actions -->
        <v-alert
          v-if="type === 'error'"
          type="warning"
          dense
          outlined
          class="mt-3"
        >
          {{ warningMessage || 'This action cannot be undone.' }}
        </v-alert>

        <!-- Confirmation input for very destructive actions -->
        <v-text-field
          v-if="requireConfirmation"
          v-model="confirmationText"
          :label="confirmationLabel"
          :placeholder="confirmationPlaceholder"
          outlined
          dense
          class="mt-3"
        />
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn
          text
          @click="$emit('input', false)"
          aria-label="Cancel"
        >
          {{ cancelText || 'Cancel' }}
        </v-btn>
        <v-btn
          :color="type"
          :disabled="requireConfirmation && confirmationText !== confirmationPlaceholder"
          @click="confirm"
          :aria-label="confirmText"
        >
          {{ confirmText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: 'ConfirmActionDialog',
  props: {
    value: Boolean,
    title: {
      type: String,
      required: true,
    },
    message: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      default: 'primary',
      validator: (v) => ['primary', 'warning', 'error'].includes(v),
    },
    confirmText: {
      type: String,
      default: 'Confirm',
    },
    cancelText: String,
    warningMessage: String,
    requireConfirmation: Boolean,
    confirmationLabel: String,
    confirmationPlaceholder: String,
  },
  data() {
    return {
      confirmationText: '',
    };
  },
  computed: {
    icon() {
      const icons = {
        primary: 'mdi-information',
        warning: 'mdi-alert',
        error: 'mdi-alert-circle',
      };
      return icons[this.type];
    },
  },
  watch: {
    value(newVal) {
      if (!newVal) {
        // Reset confirmation on close
        this.confirmationText = '';
      }
    },
  },
  methods: {
    confirm() {
      this.$emit('confirm');
      this.$emit('input', false);
    },
  },
};
</script>
