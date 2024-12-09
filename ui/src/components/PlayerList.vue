<template>
  <div class="space-y-6">
    <!-- Add Player Form - Only visible when authenticated -->
    <form v-if="isAuthenticated"
          @submit.prevent="$emit('add-player', newPlayerName)" 
          class="bg-gradient-to-r from-blue-500 to-indigo-500 p-6 rounded-xl shadow-lg space-y-4">
      <h3 class="text-xl font-bold text-white mb-4">👥 Add New Player</h3>
      <div class="flex space-x-4">
        <input
          v-model="newPlayerName"
          type="text"
          placeholder="Player Name"
          class="modern-input flex-1 bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
          required
        />
        <button type="submit" class="modern-button bg-gradient-to-r from-green-500 to-emerald-500 text-white font-bold">
          ➕ Add Player
        </button>
      </div>
    </form>

    <!-- Players List - Always visible -->
    <div v-if="!sortedPlayers.length" class="text-center text-gray-500 dark:text-gray-400">
      No players yet
    </div>
    <div v-else class="grid gap-4 sm:grid-cols-2 md:grid-cols-3">
      <div v-for="player in sortedPlayers" 
           :key="player.id" 
           class="bg-white dark:bg-gray-800 rounded-xl shadow-md hover:shadow-lg transition-shadow p-4">
        <div class="flex items-center justify-center space-x-3">
          <span class="text-2xl">👤</span>
          <span class="text-lg font-semibold text-gray-900 dark:text-white">{{ player.name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Player } from '../types';

const props = defineProps<{
  players: Player[];
  isAuthenticated: boolean;
}>();

const newPlayerName = ref('');

defineEmits<{
  (e: 'add-player', name: string): void;
}>();

const sortedPlayers = computed(() => {
  if (!props.players) return [];
  return [...props.players].sort((a, b) => a.name.localeCompare(b.name));
});
</script>
