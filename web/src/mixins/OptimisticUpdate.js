/**
 * Optimistic UI Updates Mixin
 * Allows UI to update immediately while API call is in progress
 */
import { toast } from '@/lib/toast';

export default {
  methods: {
    async optimisticUpdate(updateFn, apiCall, rollbackFn) {
      // Apply optimistic update
      const previousState = updateFn();

      try {
        // Make API call
        await apiCall();
      } catch (error) {
        // Rollback on error
        if (rollbackFn) {
          rollbackFn(previousState);
        } else {
          // Default rollback
          this.$forceUpdate();
        }

        // Show error
        toast.apiError(error);
        throw error;
      }
    },
  },
};
