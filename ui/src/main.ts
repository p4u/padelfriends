import { createApp } from 'vue';
import { createPinia } from 'pinia';
import { router } from './router';
import App from './App.vue';
import { i18n } from './i18n';
import './style.css';

// Create app instance
const app = createApp(App);

// Create pinia instance
const pinia = createPinia();

// Install plugins
app.use(pinia);
app.use(router);
app.use(i18n);

// Mount app
app.mount('#app');

// Register service worker
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/service-worker.js').then(registration => {
      console.log('Service Worker registered with scope:', registration.scope);
    }).catch(error => {
      console.error('Service Worker registration failed:', error);
    });
  });
}

// Debug i18n setup
const composer = i18n.global;
console.log('i18n debug:', {
  currentLocale: composer.locale.value,
  availableLocales: composer.availableLocales,
  loadedMessages: Object.keys(composer.messages.value)
});

// Handle errors globally
app.config.errorHandler = (err, instance, info) => {
  console.error('Global error:', err);
  console.error('Component:', instance);
  console.error('Info:', info);
};

// Handle warnings globally
app.config.warnHandler = (msg, instance, trace) => {
  console.warn('Global warning:', msg);
  console.warn('Component:', instance);
  console.warn('Trace:', trace);
};

// Handle unhandled promise rejections
window.addEventListener('unhandledrejection', event => {
  console.error('Unhandled promise rejection:', event.reason);
});

// Handle uncaught errors
window.addEventListener('error', event => {
  console.error('Uncaught error:', event.error);
});
