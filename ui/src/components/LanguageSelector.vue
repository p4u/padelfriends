<template>
  <div class="relative">
    <!-- Language Button -->
    <button 
      @click="isOpen = !isOpen"
      class="flex items-center space-x-2 px-3 py-2 rounded-lg bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
    >
      <span class="text-lg">{{ languageEmoji }}</span>
      <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ t(`languages.${currentLocale}`) }}
      </span>
    </button>

    <!-- Language Dropdown -->
    <div 
      v-if="isOpen"
      class="absolute right-0 mt-2 w-48 rounded-lg bg-white dark:bg-gray-800 shadow-lg border border-gray-200 dark:border-gray-700 py-1 z-50"
    >
      <button
        v-for="lang in languages"
        :key="lang.code"
        @click="handleLanguageSelect(lang.code)"
        class="w-full px-4 py-2 text-left hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors flex items-center space-x-2"
        :class="{ 'bg-blue-50 dark:bg-blue-900/20': lang.code === currentLocale }"
      >
        <span class="text-lg">{{ getLanguageEmoji(lang.code) }}</span>
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
          {{ t(`languages.${lang.code}`) }}
        </span>
      </button>
    </div>

    <!-- Backdrop -->
    <div 
      v-if="isOpen"
      class="fixed inset-0 z-40"
      @click="isOpen = false"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue';
import { useI18n, type Language } from '../i18n';

const { t, locale, setLocale } = useI18n();
const isOpen = ref(false);

const languages = [
  { code: 'en' as Language },
  { code: 'es' as Language },
  { code: 'ca' as Language }
] as const;

const currentLocale = computed(() => locale.value);

const getLanguageEmoji = (code: Language): string => {
  switch (code) {
    case 'en': return '🇬🇧';
    case 'es': return '🇪🇸';
    case 'ca': return '🏴󠁥󠁳󠁣󠁴󠁿';
    default: return '🌐';
  }
};

const languageEmoji = computed(() => getLanguageEmoji(currentLocale.value as Language));

const handleLanguageSelect = async (code: Language) => {
  console.log('Changing language to:', code);
  console.log('Current locale before change:', currentLocale.value);
  
  setLocale(code);
  isOpen.value = false;
  
  await nextTick();
  console.log('Current locale after change:', currentLocale.value);
  console.log('Current translation:', t('app.name'));
};
</script>
