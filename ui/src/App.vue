<template>
  <div class="min-h-screen bg-gradient-to-b from-background to-gray-50 dark:from-background dark:to-gray-900 flex flex-col">
    <TopMenu v-if="showTopMenu" />
    <div :class="[
      'max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex-grow',
      { 'mt-16': showTopMenu }
    ]">
      <router-view />
    </div>
    <footer class="w-full py-4 text-center text-sm text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-900">
      <a href="https://github.com/p4u/padelfriends" target="_blank" rel="noopener noreferrer" class="hover:text-gray-900 dark:hover:text-gray-200">
        This application is opensource - View on GitHub
      </a>
    </footer>
    <InstallButton />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { InstallButton, TopMenu } from './components';

const route = useRoute();

const showTopMenu = computed(() => {
  return route.path.startsWith('/group/');
});
</script>
