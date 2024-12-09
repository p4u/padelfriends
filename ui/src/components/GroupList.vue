<template>
  <div class="space-y-2">
    <div v-if="loading" class="text-center text-gray-900 dark:text-white">
      Loading groups... 🔄
    </div>
    <div v-else-if="error" class="text-center text-red-500">
      {{ error }} 😢
    </div>
    <div v-else-if="groups.length === 0" class="text-center text-gray-900 dark:text-white">
      No groups available 😢
    </div>
    <div v-else class="space-y-2">
      <div
        v-for="group in groups"
        :key="group?.id || 'unknown'"
        class="modern-container bg-white dark:bg-gray-800 p-4 cursor-pointer hover:border-primary transition-colors"
        @click="handleGroupSelect(group)"
      >
        <span class="text-gray-900 dark:text-white">{{ group?.name || 'Unnamed Group' }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import type { Group } from '../types';

const props = defineProps<{
  groups: Group[];
  loading: boolean;
  error: string | null;
}>();

const emit = defineEmits<{
  (e: 'select', group: Group): void;
}>();

const handleGroupSelect = (group: Group) => {
  if (!group || !group.name) {
    console.error('Invalid group data:', group);
    return;
  }
  emit('select', group);
};

onMounted(() => {
  console.log('GroupList mounted with groups:', props.groups);
});
</script>
