<template>
  <div class="space-y-4">
    <div v-if="!savedGroups.length" class="text-center text-gray-500 dark:text-gray-400">
      {{ t('groups.noSavedGroups') }}
    </div>
    <div v-else class="grid gap-4 grid-cols-1 sm:grid-cols-2">
      <div 
        v-for="group in savedGroups" 
        :key="group.name"
        class="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg space-y-4"
      >
        <div class="flex justify-between items-center">
          <h3 class="font-bold text-gray-900 dark:text-white">
            {{ group.name }}
          </h3>
          <button 
            @click="joinGroup(group)"
            class="modern-button bg-gradient-to-r from-blue-500 to-indigo-500 text-white"
          >
            {{ t('groups.join') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from '../i18n';

interface SavedGroup {
  name: string;
  password: string;
}

const emit = defineEmits<{
  (e: 'join-group', name: string, password: string): void;
}>();

const { t } = useI18n();

const savedGroups = computed<SavedGroup[]>(() => {
  try {
    const savedGroupsStr = localStorage.getItem('savedGroups');
    return savedGroupsStr ? JSON.parse(savedGroupsStr) : [];
  } catch (error) {
    console.error('Error parsing saved groups:', error);
    return [];
  }
});

const joinGroup = (group: SavedGroup) => {
  emit('join-group', group.name, group.password);
};
</script>
