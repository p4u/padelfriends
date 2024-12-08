<template>
  <div class="space-y-4">
    <form @submit.prevent="$emit('add-player', newPlayerName)" class="flex space-x-4">
      <input
        v-model="newPlayerName"
        type="text"
        placeholder="Player Name"
        class="retro-input flex-1"
        required
      />
      <button type="submit" class="retro-button">Add Player</button>
    </form>
    <ul class="space-y-2">
      <li v-for="player in players" :key="player.id" class="p-2 bg-gray-100 rounded">
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