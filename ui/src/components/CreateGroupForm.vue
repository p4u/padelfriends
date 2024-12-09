<template>
  <form @submit.prevent="handleSubmit" class="modern-container bg-white dark:bg-gray-800">
    <h2 class="text-xl font-bold mb-4 text-gray-900 dark:text-white">
      {{ t('groups.createNew') }}
    </h2>
    <div class="space-y-4">
      <!-- Group Name -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          {{ t('groups.name') }}
        </label>
        <input 
          v-model="groupName"
          type="text"
          required
          class="modern-input w-full text-gray-900 dark:text-white bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
          :placeholder="t('groups.name')"
        />
      </div>

      <!-- Password -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          {{ t('auth.password') }}
        </label>
        <input 
          v-model="password"
          type="password"
          required
          class="modern-input w-full text-gray-900 dark:text-white bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
          :placeholder="t('auth.password')"
        />
      </div>

      <!-- Submit Button -->
      <button 
        type="submit"
        class="modern-button w-full bg-gradient-to-r from-blue-500 to-indigo-500 text-white"
      >
        {{ t('groups.create') }}
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from '../i18n';

const emit = defineEmits<{
  (e: 'create-group', name: string, password: string): void;
}>();

const { t } = useI18n();
const groupName = ref('');
const password = ref('');

const handleSubmit = () => {
  if (groupName.value.trim() && password.value.trim()) {
    emit('create-group', groupName.value.trim(), password.value.trim());
    groupName.value = '';
    password.value = '';
  }
};
</script>
