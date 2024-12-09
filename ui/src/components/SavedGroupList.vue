<template>
  <div v-if="groups.length > 0" class="space-y-2">
    <div 
      v-for="group in sortedGroups" 
      :key="group.id"
      class="modern-container bg-white dark:bg-gray-800 p-4 cursor-pointer hover:border-primary transition-colors"
      @click="$emit('select', group)"
    >
      <div class="flex justify-between items-center">
        <span class="text-gray-900 dark:text-white font-medium">{{ group.name }}</span>
        <button 
          class="text-red-500 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
          @click.stop="$emit('remove', group.id)"
        >
          ❌
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface SavedGroup {
  id: string;
  name: string;
  password: string;
  lastAccessed: string;
}

const props = defineProps<{
  groups: SavedGroup[];
}>();

defineEmits<{
  (e: 'select', group: SavedGroup): void;
  (e: 'remove', id: string): void;
}>();

const sortedGroups = computed(() => 
  [...props.groups].sort((a, b) => 
    new Date(b.lastAccessed).getTime() - new Date(a.lastAccessed).getTime()
  )
);
</script>
