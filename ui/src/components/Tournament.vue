<template>
  <div>
    <!-- Tournament Button -->
    <button 
      @click="showModal = true"
      class="w-full mt-8 p-4 bg-gradient-to-r from-purple-500 to-pink-500 text-white font-bold rounded-xl shadow-lg hover:shadow-xl transition-shadow"
    >
      {{ t('matches.startTournament') }}
    </button>

    <!-- Modal -->
    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto p-6">
        <div class="space-y-6">
          <!-- Header -->
          <div class="flex justify-between items-center">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white">
              {{ t('tournament.setup') }}
            </h3>
            <button @click="showModal = false" class="text-gray-500 hover:text-gray-700 dark:text-gray-400">
              ❌
            </button>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="text-red-500 text-sm">
            {{ error }}
          </div>

          <!-- Player Selection -->
          <div class="space-y-4">
            <h4 class="font-medium text-gray-900 dark:text-white">
              {{ t('tournament.selectPlayers') }}
            </h4>
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
              <label 
                v-for="player in players" 
                :key="player.id"
                class="flex items-center space-x-2 p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 cursor-pointer"
              >
                <input 
                  type="checkbox" 
                  :value="player.id"
                  v-model="selectedPlayerIds"
                  class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                >
                <span class="text-gray-900 dark:text-white">{{ player.name }}</span>
              </label>
            </div>
            <div class="text-sm text-gray-500">
              {{ t('tournament.selectedPlayers') }}: {{ selectedPlayerIds.length }}
              <span v-if="!isValidPlayerCount" class="text-red-500">
                ({{ t('tournament.playerMultiple') }})
              </span>
            </div>
          </div>

          <!-- Number of Matches -->
          <div class="space-y-2">
            <label class="block font-medium text-gray-900 dark:text-white">
              {{ t('tournament.numberOfMatches') }}
            </label>
            <input 
              type="number" 
              v-model.number="numberOfMatches"
              min="1"
              class="modern-input w-full"
              :disabled="!isValidPlayerCount"
            >
          </div>

          <!-- Actions -->
          <div class="flex justify-end space-x-3">
            <button 
              @click="showModal = false"
              class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 dark:text-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600"
            >
              {{ t('common.cancel') }}
            </button>
            <button 
              @click="generateMatches"
              :disabled="!isValid"
              class="px-4 py-2 text-white bg-gradient-to-r from-purple-500 to-pink-500 rounded-lg hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ t('tournament.generateMatches') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from '../i18n';
import type { Player } from '../types';

const props = defineProps<{
  players: Player[];
}>();

const emit = defineEmits<{
  (e: 'create-matches', matches: Array<string[]>): void;
}>();

const { t } = useI18n();
const showModal = ref(false);
const selectedPlayerIds = ref<string[]>([]);
const numberOfMatches = ref(1);
const error = ref('');

const isValidPlayerCount = computed(() => {
  return selectedPlayerIds.value.length >= 4 && selectedPlayerIds.value.length % 4 === 0;
});

const isValid = computed(() => {
  return isValidPlayerCount.value && numberOfMatches.value > 0;
});

interface PlayerPairings {
  [key: string]: Set<string>;
}

interface PlayerCounts {
  [key: string]: number;
}

const getPairingScore = (
  match: string[], 
  existingPairings: PlayerPairings,
  playerCounts: PlayerCounts
): number => {
  let score = 0;
  
  for (let i = 0; i < match.length; i++) {
    for (let j = i + 1; j < match.length; j++) {
      if (existingPairings[match[i]]?.has(match[j])) {
        score += 10;
      }
    }
  }

  const counts = match.map(id => playerCounts[id] || 0);
  const maxCount = Math.max(...counts);
  const minCount = Math.min(...counts);
  score += (maxCount - minCount) * 5;

  return score;
};

const updatePairings = (match: string[], pairings: PlayerPairings, counts: PlayerCounts) => {
  for (let i = 0; i < match.length; i++) {
    if (!pairings[match[i]]) {
      pairings[match[i]] = new Set();
    }
    counts[match[i]] = (counts[match[i]] || 0) + 1;
    
    for (let j = 0; j < match.length; j++) {
      if (i !== j) {
        pairings[match[i]].add(match[j]);
      }
    }
  }
};

const generateOptimalMatches = (playerIds: string[], count: number): Array<string[]> => {
  const result: Array<string[]> = [];
  const pairings: PlayerPairings = {};
  const playerCounts: PlayerCounts = {};

  playerIds.forEach(id => {
    pairings[id] = new Set();
    playerCounts[id] = 0;
  });

  for (let matchIndex = 0; matchIndex < count; matchIndex++) {
    let bestMatch: string[] | null = null;
    let bestScore = Infinity;

    const attempts = 1000;
    for (let attempt = 0; attempt < attempts; attempt++) {
      const availablePlayers = [...playerIds];
      const match: string[] = [];

      while (match.length < 4 && availablePlayers.length > 0) {
        const randomIndex = Math.floor(Math.random() * availablePlayers.length);
        match.push(availablePlayers[randomIndex]);
        availablePlayers.splice(randomIndex, 1);
      }

      const score = getPairingScore(match, pairings, playerCounts);
      if (score < bestScore) {
        bestScore = score;
        bestMatch = match;
      }

      if (bestScore === 0) break;
    }

    if (bestMatch) {
      result.push(bestMatch);
      updatePairings(bestMatch, pairings, playerCounts);
    }
  }

  return result;
};

const generateMatches = () => {
  if (!isValid.value) return;

  try {
    const matches = generateOptimalMatches(selectedPlayerIds.value, numberOfMatches.value);
    emit('create-matches', matches);
    showModal.value = false;
    selectedPlayerIds.value = [];
    numberOfMatches.value = 1;
    error.value = '';
  } catch (err) {
    error.value = t('errors.failedToCreate');
  }
};
</script>
