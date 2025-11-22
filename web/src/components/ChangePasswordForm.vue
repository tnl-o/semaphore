<template>
  <v-form
    ref="form"
    lazy-validation
    v-model="formValid"
  >
    <v-alert
      :value="formError"
      color="error"
      class="pb-2"
    >{{ formError }}</v-alert>

    <SecurePasswordInput
      v-model="item.password"
      :label="$t('password2')"
      :rules="[v => !!v || $t('password_required')]"
      required
      :disabled="formSaving"
    />
  </v-form>
</template>
<script>
import ItemFormBase from '@/components/ItemFormBase';
import SecurePasswordInput from '@/components/SecurePasswordInput.vue';

export default {
  mixins: [ItemFormBase],
  components: {
    SecurePasswordInput,
  },
  methods: {
    async loadData() {
      this.item = {};
    },

    getItemsUrl() {
      return null;
    },

    getSingleItemUrl() {
      return null;
    },

    getRequestOptions() {
      return {
        method: 'post',
        url: `/api/users/${this.itemId}/password`,
      };
    },
  },
};
</script>
