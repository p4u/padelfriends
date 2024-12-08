<template>
  <div class="space-y-4">
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div class="grid grid-cols-2 gap-4">
        <div v-for="(team, index) in ['Team 1', 'Team 2']" :key="team">
          <h3 class="mb-2">{{ team }}</h3>
          <div class="space-y-2">
            <select
              v-model="selectedPlayers[index][0]"
              class="retro-input w-full"
              required
            >
              <option v-for="player in availablePlayers" :key="player.id" :value="player.id">
                {{ player.name }}
              </option>
            </select>
            <select
              v-model="selectedPlayers[index][1]"
              class="retro-input w-full"
              required
            >
              <option v-for="player in availablePlayers" :key="player.id" :value="player.id">
                {{ player.name }}
              </option>
            </select>
          </div>
        </div>
      </div>
      <button type="submit" class="retro-button w-full">Create Match</button>
    </form>

    <div class="space-y-4">
      <div v-for="match in matches" :key="match.id" class="p-4 bg-gray-100 rounded">
        <div class="flex justify-between items-center">
          <div>
            Team 1: {{ match.players ? getPlayerNames(match.players.slice(0, 2)) : 'N/A' }}
            vs
            Team 2: {{ match.players ? getPlayerNames(match.players.slice(2, 4)) : 'N/A' }}
          </div>
          <div v-if="match.status === 'completed'">
            Score: {{ match.score_team1 }} - {{ match.score_team2 }}
          </div>
          <button
            v-else
            @click="$emit('submit-score', match)"
            class="retro-button"
          >
            Submit Score
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Match, Player } from '../types';

const props = defineProps<{
  matches: Match[];
  players: Player[];
}>();

const emit = defineEmits<{
  (e: 'create-match', playerIds: string[]): void;
  (e: 'submit-score', match: Match): void;
}>();

const selectedPlayers = ref([['', ''], ['', '']]);

const availablePlayers = computed(() => props.players);

const getPlayerNames = (playerIds: string[]) => {
  return playerIds
    .map(id => props.players.find(p => p.id === id)?.name)
    .filter(Boolean)
    .join(' & ');
};

const handleSubmit = () => {
  const playerIds = selectedPlayers.value.flat();
  emit('create-match', playerIds);
  selectedPlayers.value = [['', ''], ['', '']];
};
</script>
