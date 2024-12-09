<template>
  <button
    v-if="showInstallButton"
    @click="installPWA"
    class="fixed bottom-4 right-4 px-4 py-2 bg-gradient-to-r from-blue-500 to-indigo-500 text-white rounded-lg shadow-lg flex items-center space-x-2 hover:shadow-xl transition-shadow"
  >
    <span class="text-xl">📱</span>
    <span>Install App</span>
  </button>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

const deferredPrompt = ref<any>(null);
const showInstallButton = ref(false);

const handleBeforeInstallPrompt = (e: Event) => {
  e.preventDefault();
  deferredPrompt.value = e;
  showInstallButton.value = true;
};

const installPWA = async () => {
  if (!deferredPrompt.value) return;
  
  deferredPrompt.value.prompt();
  const { outcome } = await deferredPrompt.value.userChoice;
  
  if (outcome === 'accepted') {
    showInstallButton.value = false;
  }
  
  deferredPrompt.value = null;
};

onMounted(() => {
  window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt);
});

onUnmounted(() => {
  window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt);
});
</script>
