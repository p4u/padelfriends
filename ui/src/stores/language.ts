import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useI18n, type Language } from '../i18n';

export const useLanguageStore = defineStore('language', () => {
  const { locale } = useI18n();

  const availableLanguages = [
    { code: 'en' as Language, name: 'English' },
    { code: 'es' as Language, name: 'Español' },
    { code: 'ca' as Language, name: 'Català' }
  ] as const;

  return {
    locale,
    availableLanguages
  };
});
