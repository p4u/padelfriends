<template>
  <div class="space-y-4">
    <form @submit.prevent="$emit('add-player', newPlayerName)" class="flex space-x-4">
      <input
        v-model="newPlayerName"
        type="text"
        placeholder="Player Name"
        class="modern-input flex-1"
        required
      />
      <button type="submit" class="modern-button">Add Player</button>
    </form>
    <ul class="space-y-2">
      <li v-for="player in players" :key="player.id" class="p-3 bg-white dark:bg-gray-800 rounded-lg shadow-sm text-center text-gray-900 dark:text-white font-medium">
        {{ player.name }}
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { Player } from '../types';

defineProps<{
  players: Player[];
}>();

const newPlayerName = ref('');

defineEmits<{
  (e: 'add-player', name: string): void;
}>();
</script>
