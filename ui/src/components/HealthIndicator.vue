<template>
  <div class="fixed bottom-4 right-4 flex items-center space-x-2 px-3 py-2 rounded-lg bg-white dark:bg-gray-800 shadow-lg border border-gray-200 dark:border-gray-700">
    <div 
      class="w-3 h-3 rounded-full"
      :class="{
        'bg-green-500': isHealthy,
        'bg-red-500': !isHealthy,
        'animate-pulse': !isHealthy
      }"
    ></div>
    <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
      {{ t(isHealthy ? 'common.connected' : 'common.disconnected') }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useI18n } from '../i18n';
import { healthApi } from '../api';

const { t } = useI18n();
const isHealthy = ref(true);
let healthCheckInterval: number | undefined;

const checkHealth = async () => {
  try {
    await healthApi.check();
    isHealthy.value = true;
  } catch (error) {
    isHealthy.value = false;
  }
};

onMounted(() => {
  checkHealth();
  healthCheckInterval = window.setInterval(checkHealth, 30000);
});

onUnmounted(() => {
  if (healthCheckInterval) {
    clearInterval(healthCheckInterval);
  }
});
</script>
