<template>
  <div class="fixed top-0 left-0 right-0 bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm shadow-md z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center space-x-4">
          <img src="/icons/icon-192x192.png" alt="Padel Friends Logo" class="h-10 w-10" />
          <span class="text-xl font-bold text-gray-900 dark:text-white">padelfriends.xyz</span>
        </div>
        <div class="flex items-center space-x-2">
          <button 
            @click="refreshPage"
            class="px-4 py-2 rounded-lg bg-blue-500/10 hover:bg-blue-500/20 text-blue-600 dark:text-blue-400 transition-colors"
          >
            🔄 Refresh
          </button>
          <router-link 
            to="/" 
            class="px-4 py-2 rounded-lg bg-blue-500/10 hover:bg-blue-500/20 text-blue-600 dark:text-blue-400 transition-colors"
          >
            Home
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useGroupStore } from '../stores/group';

const router = useRouter();
const groupStore = useGroupStore();

const refreshPage = async () => {
  if (router.currentRoute.value.path.startsWith('/group/')) {
    // Reload group data
    await Promise.all([
      groupStore.loadPlayers(),
      groupStore.loadMatches(),
      groupStore.loadStatistics()
    ]);
  } else {
    // For other pages, just reload the current route
    router.go(0);
  }
};
</script>
