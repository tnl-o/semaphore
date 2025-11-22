<template>
  <div>
    <v-skeleton-loader
      v-if="loading"
      :type="type"
      :tile="tile"
      :height="height"
    />
    <slot v-else />
  </div>
</template>

<script>
export default {
  name: 'ContentSkeleton',
  props: {
    loading: Boolean,
    type: {
      type: String,
      default: 'card',
    },
    tile: Boolean,
    height: [String, Number],
  },
};
</script>
