<template>
  <div class="flex items-center space-x-2">
    <div 
      class="w-3 h-3 rounded-full"
      :class="[
        isHealthy ? 'bg-green-500 animate-pulse' : 'bg-red-500',
        'shadow-lg',
        'border border-black'
      ]"
    ></div>
    <span class="text-xs">
      {{ isHealthy ? '🟢 Connected' : '🔴 Disconnected' }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { healthApi } from '../api';

const isHealthy = ref(false);
let healthCheckInterval: number;

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
  healthCheckInterval = setInterval(checkHealth, 30000) as unknown as number;
});

onUnmounted(() => {
  clearInterval(healthCheckInterval);
});
</script>