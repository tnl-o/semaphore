<template>
  <div class="task-log-records__record">
    <div class="task-log-records__time">
      {{ source.time | formatTime }}
    </div>
    <SafeLogOutput :output="source.output" />
  </div>
</template>
<script>
import SafeLogOutput from '@/components/SafeLogOutput.vue';

export default {
  components: {
    SafeLogOutput,
  },
  props: {
    index: { // index of current item
      type: Number,
    },
    source: { // here is: {uid: 'unique_1', text: 'abc'}
      type: Object,
      default() {
        return {};
      },
    },
  },
};
</script>
