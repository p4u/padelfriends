<template>
  <div class="space-y-4">
    <div v-if="!groups.length" class="text-center text-gray-500 dark:text-gray-400">
      {{ t('groups.noGroups') }}
    </div>
    <div v-else class="grid gap-4 grid-cols-1 sm:grid-cols-2">
      <div 
        v-for="group in groups" 
        :key="group.name"
        class="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg space-y-4"
      >
        <div class="flex justify-between items-center">
          <h3 class="font-bold text-gray-900 dark:text-white">
            {{ group.name }}
          </h3>
          <span class="text-sm text-gray-500 dark:text-gray-400">
            {{ formatDate(group.created_at) }}
          </span>
        </div>
        
        <div class="flex items-center space-x-2">
          <input 
            v-model="passwords[group.name]"
            type="password"
            :placeholder="t('auth.password')"
            class="modern-input flex-1 text-gray-900 dark:text-white bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
          />
          <button 
            @click="joinGroup(group.name)"
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
import { ref } from 'vue';
import { useI18n } from '../i18n';
import type { Group } from '../types';

const props = defineProps<{
  groups: Group[];
}>();

const emit = defineEmits<{
  (e: 'join-group', name: string, password: string): void;
}>();

const { t } = useI18n();
const passwords = ref<{ [key: string]: string }>({});

const formatDate = (timestamp: string) => {
  const date = new Date(timestamp);
  const now = new Date();
  const diffTime = Math.abs(now.getTime() - date.getTime());
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

  if (diffDays === 0) {
    return t('common.today') + ', ' + date.toLocaleTimeString();
  } else if (diffDays === 1) {
    return t('common.yesterday') + ', ' + date.toLocaleTimeString();
  } else {
    return date.toLocaleDateString();
  }
};

const joinGroup = (name: string) => {
  const password = passwords.value[name];
  if (!password) {
    alert(t('errors.invalidPassword'));
    return;
  }
  emit('join-group', name, password);
  passwords.value[name] = '';
};
</script>
