import { createI18n } from 'vue-i18n';
import type { I18n, Composer } from 'vue-i18n';
import en from '../locales/en.json';
import es from '../locales/es.json';
import ca from '../locales/ca.json';

export type Language = 'en' | 'es' | 'ca';
type MessageSchema = typeof en;

const storedLang = localStorage.getItem('language');
const browserLang = navigator.language.split('-')[0];
const defaultLocale = (storedLang && ['en', 'es', 'ca'].includes(storedLang) 
  ? storedLang 
  : ['en', 'es', 'ca'].includes(browserLang) 
    ? browserLang 
    : 'en') as Language;

// Create i18n instance
export const i18n = createI18n({
  legacy: false,
  locale: defaultLocale,
  fallbackLocale: 'en',
  messages: {
    en,
    es,
    ca
  },
  globalInjection: true,
  silentTranslationWarn: true,
  missingWarn: false,
  fallbackWarn: false,
  sync: true
});

// Type-safe composable
export function useI18n() {
  const composer = i18n.global;

  const setLocale = (newLocale: Language) => {
    try {
      // @ts-ignore - Vue I18n types are not perfect
      composer.locale.value = newLocale;
      document.documentElement.lang = newLocale;
      localStorage.setItem('language', newLocale);
      
      // Force reactivity update
      window.dispatchEvent(new Event('languageChanged'));
      
      // Log for debugging
      console.log('Language changed:', {
        newLocale,
        // @ts-ignore - Vue I18n types are not perfect
        currentLocale: composer.locale.value,
        availableLocales: composer.availableLocales,
        storedLocale: localStorage.getItem('language'),
        // @ts-ignore - Vue I18n types are not perfect
        messages: Object.keys(composer.messages.value)
      });
    } catch (error) {
      console.error('Error changing language:', error);
    }
  };

  return {
    t: composer.t,
    // @ts-ignore - Vue I18n types are not perfect
    locale: composer.locale,
    setLocale
  };
}

export function setupI18n() {
  const composer = i18n.global;
  
  // Set initial HTML lang attribute
  document.documentElement.lang = defaultLocale;
  
  // Log initial setup
  console.log('i18n setup:', {
    defaultLocale,
    storedLang,
    browserLang,
    // @ts-ignore - Vue I18n types are not perfect
    currentLocale: composer.locale.value,
    availableLocales: composer.availableLocales,
    // @ts-ignore - Vue I18n types are not perfect
    messages: Object.keys(composer.messages.value)
  });
  
  return i18n;
}
