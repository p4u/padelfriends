<template>
  <div class="space-y-8">
    <!-- Add Player Form -->
    <form @submit.prevent="handleSubmit" class="space-y-4 bg-gradient-to-r from-blue-500 to-indigo-500 p-6 rounded-xl shadow-lg">
      <h3 class="text-xl font-bold text-white mb-4">{{ t('players.addPlayer') }}</h3>
      <div class="space-y-4">
        <input 
          v-model="newPlayerName"
          type="text"
          :placeholder="t('players.playerName')"
          class="modern-input w-full text-gray-900 dark:text-white bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
          required
        />
        <button 
          type="submit" 
          class="modern-button w-full bg-gradient-to-r from-green-500 to-emerald-500 text-white font-bold"
        >
          {{ t('players.addPlayer') }}
        </button>
      </div>
    </form>

    <!-- Players List -->
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
          {{ t('players.playerList') }}
        </h2>
      </div>
      
      <div class="p-6">
        <div v-if="!players.length" class="text-center text-gray-500 dark:text-gray-400">
          {{ t('players.noPlayers') }}
        </div>
        <div v-else class="grid gap-4 grid-cols-1 sm:grid-cols-2 md:grid-cols-3">
          <div 
            v-for="player in players" 
            :key="player.id"
            class="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg flex items-center justify-between"
          >
            <span class="font-medium text-gray-900 dark:text-white">
              {{ player.name }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from '../i18n';
import type { Player } from '../types';

const props = defineProps<{
  players: Player[];
}>();

const emit = defineEmits<{
  (e: 'add-player', name: string): void;
}>();

const { t } = useI18n();
const newPlayerName = ref('');

const handleSubmit = () => {
  if (newPlayerName.value.trim()) {
    emit('add-player', newPlayerName.value.trim());
    newPlayerName.value = '';
  }
};
</script>
