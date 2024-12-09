<template>
  <div class="space-y-6">
    <!-- Create Match Form -->
    <form @submit.prevent="handleSubmit" class="space-y-4 bg-gradient-to-r from-blue-500 to-indigo-500 p-6 rounded-xl shadow-lg">
      <h3 class="text-xl font-bold text-white mb-4">🎾 New Match</h3>
      <div class="grid grid-cols-2 gap-6">
        <div v-for="(team, index) in ['🏆', '🎯']" :key="team" class="space-y-3">
          <div class="flex items-center space-x-2">
            <span class="text-2xl">{{ team }}</span>
            <span class="text-white font-medium">Team {{ index + 1 }}</span>
          </div>
          <div class="space-y-2">
            <select
              v-model="selectedPlayers[index][0]"
              class="modern-input w-full text-gray-900 dark:text-white bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
              required
              @change="validateSelection"
            >
              <option value="">Select Player 1</option>
              <option 
                v-for="player in availablePlayersForSelect(index, 0)" 
                :key="player.id" 
                :value="player.id"
              >
                {{ player.name }}
              </option>
            </select>
            <select
              v-model="selectedPlayers[index][1]"
              class="modern-input w-full text-gray-900 dark:text-white bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
              required
              @change="validateSelection"
            >
              <option value="">Select Player 2</option>
              <option 
                v-for="player in availablePlayersForSelect(index, 1)" 
                :key="player.id" 
                :value="player.id"
              >
                {{ player.name }}
              </option>
            </select>
          </div>
        </div>
      </div>
      <div v-if="selectionError" class="text-red-200 text-sm mt-2">
        {{ selectionError }}
      </div>
      <button 
        type="submit" 
        class="modern-button w-full bg-gradient-to-r from-green-500 to-emerald-500 text-white font-bold"
        :disabled="!!selectionError"
      >
        🎮 Start Match
      </button>
    </form>

    <!-- Matches List -->
    <div class="space-y-4">
      <div v-for="match in sortedMatches" :key="match.id" 
           class="p-6 bg-white dark:bg-gray-800 rounded-xl shadow-md hover:shadow-lg transition-shadow">
        <div class="space-y-4">
          <!-- Match Header -->
          <div class="text-sm text-gray-500 dark:text-gray-400">
            {{ formatDate(match.timestamp) }}
          </div>

          <!-- Match Content -->
          <div class="flex items-center justify-between">
            <!-- Teams -->
            <div class="flex-1 grid grid-cols-[1fr,auto,1fr] items-center gap-4">
              <!-- Team 1 -->
              <div class="text-right">
                <div class="font-bold text-blue-600 dark:text-blue-400">
                  {{ getTeamNames(match.team1) }}
                </div>
              </div>

              <!-- VS -->
              <div class="px-4 py-2 rounded-full bg-gray-100 dark:bg-gray-700 font-bold text-gray-600 dark:text-gray-300">
                VS
              </div>

              <!-- Team 2 -->
              <div class="text-left">
                <div class="font-bold text-indigo-600 dark:text-indigo-400">
                  {{ getTeamNames(match.team2) }}
                </div>
              </div>
            </div>

            <!-- Score or Action -->
            <div class="ml-6">
              <div v-if="match.status === 'completed'" 
                   class="text-2xl font-bold bg-gradient-to-r from-blue-500 to-indigo-500 text-white px-4 py-2 rounded-lg">
                {{ match.score_team1 }} - {{ match.score_team2 }}
              </div>
              <button
                v-else
                @click="$emit('submit-score', match)"
                class="modern-button bg-gradient-to-r from-green-500 to-emerald-500 text-white font-bold"
              >
                🎯 Submit Score
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Match, Player, PlayerInfo } from '../types';

const props = defineProps<{
  matches: Match[];
  players: Player[];
}>();

const emit = defineEmits<{
  (e: 'create-match', playerIds: string[]): void;
  (e: 'submit-score', match: Match): void;
}>();

const selectedPlayers = ref([['', ''], ['', '']]);
const selectionError = ref('');

const availablePlayersForSelect = (teamIndex: number, playerIndex: number) => {
  return props.players.filter(player => {
    // Get all currently selected players except the current selection
    const selectedIds = selectedPlayers.value.flat().filter((id, idx) => {
      const currentTeamIndex = Math.floor(idx / 2);
      const currentPlayerIndex = idx % 2;
      return id && (teamIndex !== currentTeamIndex || playerIndex !== currentPlayerIndex);
    });
    // Player is available if not already selected
    return !selectedIds.includes(player.id);
  });
};

const validateSelection = () => {
  const selectedIds = selectedPlayers.value.flat().filter(id => id !== '');
  const uniqueIds = new Set(selectedIds);
  
  if (selectedIds.length !== uniqueIds.size) {
    selectionError.value = 'Each player can only be selected once';
  } else {
    selectionError.value = '';
  }
};

const sortedMatches = computed(() => {
  return [...props.matches].sort((a, b) => 
    new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
  );
});

const getTeamNames = (team: PlayerInfo[]) => {
  return team.map(player => player.name).join(' & ');
};

const formatDate = (timestamp: string) => {
  const date = new Date(timestamp);
  const now = new Date();
  const diffTime = Math.abs(now.getTime() - date.getTime());
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

  if (diffDays === 0) {
    return 'Today, ' + date.toLocaleTimeString('en-US', { 
      hour: '2-digit', 
      minute: '2-digit',
      hour12: true 
    });
  } else if (diffDays === 1) {
    return 'Yesterday, ' + date.toLocaleTimeString('en-US', { 
      hour: '2-digit', 
      minute: '2-digit',
      hour12: true 
    });
  } else {
    return date.toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit', 
      minute: '2-digit',
      hour12: true 
    });
  }
};

const handleSubmit = () => {
  if (selectionError.value) return;
  
  const playerIds = selectedPlayers.value.flat();
  emit('create-match', playerIds);
  selectedPlayers.value = [['', ''], ['', '']];
  selectionError.value = '';
};
</script>
