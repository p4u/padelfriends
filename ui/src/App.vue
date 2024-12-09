<template>
  <div class="min-h-screen bg-gradient-to-b from-background to-gray-50 dark:from-background dark:to-gray-900">
    <!-- Header with Language Selector -->
    <header class="fixed top-0 left-0 right-0 bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm z-40 border-b border-gray-200 dark:border-gray-800">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-2 flex justify-between items-center safe-top safe-left safe-right">
        <h1 class="text-xl font-bold text-gray-900 dark:text-white">
          {{ t('app.name') }}
        </h1>
        <LanguageSelector />
      </div>
    </header>

    <!-- Main Content -->
    <main class="pt-16 pb-8 min-h-screen">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 safe-left safe-right">
        <router-view v-slot="{ Component }">
          <component 
            :is="Component" 
            :key="`${currentLocale}-${$route.fullPath}`" 
          />
        </router-view>
      </div>
    </main>

    <!-- Footer -->
    <footer class="bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm border-t border-gray-200 dark:border-gray-800">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 safe-bottom safe-left safe-right">
        <p class="text-sm text-center text-gray-600 dark:text-gray-400">
          {{ t('app.tagline') }}
        </p>
      </div>
    </footer>

    <!-- Health Indicator -->
    <HealthIndicator />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from './i18n';
import { LanguageSelector, HealthIndicator } from './components';

const { t, locale } = useI18n();
const currentLocale = computed(() => locale.value);
</script>

<style>
#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

main {
  flex: 1;
}

/* Improved scrolling */
* {
  -webkit-overflow-scrolling: touch;
}

/* Prevent pull-to-refresh in PWA mode */
@media (display-mode: standalone) {
  body {
    overscroll-behavior-y: none;
  }
}
</style>
