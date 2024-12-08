<template>
  <div class="space-y-2">
    <div v-if="loading" class="text-center">
      Loading groups... 🔄
    </div>
    <div v-else-if="error" class="text-center text-red-500">
      {{ error }} 😢
    </div>
    <div v-else-if="groups.length === 0" class="text-center">
      No groups available 😢
    </div>
    <div v-else class="space-y-2">
      <div
        v-for="group in groups"
        :key="group.id"
        class="retro-container p-4 cursor-pointer hover:border-primary transition-colors"
        @click="$emit('select', group)"
      >
        <span>{{ group.name }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Group } from '../types';

defineProps<{
  groups: Group[];
  loading: boolean;
  error: string | null;
}>();

defineEmits<{
  (e: 'select', group: Group): void;
}>();
</script>