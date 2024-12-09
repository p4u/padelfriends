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

    <!-- Tournament Generator -->
    <Tournament 
      :players="players"
      @create-matches="handleRandomMatches"
    />

    <!-- Matches List -->
    <div class="space-y-4">
      <div v-for="match in matches" :key="match.id" 
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
              <div v-else-if="match.status === 'pending'" class="flex space-x-2">
                <button
                  @click="showSubmitScore(match)"
                  class="modern-button bg-gradient-to-r from-green-500 to-emerald-500 text-white font-bold"
                >
                  🎯 Submit Score
                </button>
                <button
                  @click="handleCancel(match)"
                  class="modern-button bg-gradient-to-r from-red-500 to-pink-500 text-white font-bold"
                >
                  ❌ Cancel
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex justify-center space-x-2 mt-6">
        <button 
          @click="changePage(currentPage - 1)"
          :disabled="currentPage === 1"
          class="px-4 py-2 rounded-lg bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 disabled:opacity-50"
        >
          Previous
        </button>
        <div class="flex items-center space-x-2">
          <span class="text-gray-600 dark:text-gray-400">
            Page {{ currentPage }} of {{ totalPages }}
          </span>
        </div>
        <button 
          @click="changePage(currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="px-4 py-2 rounded-lg bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 disabled:opacity-50"
        >
          Next
        </button>
      </div>
    </div>

    <!-- Submit Score Modal -->
    <SubmitScore
      :show="showScoreModal"
      :match="selectedMatch"
      @close="closeSubmitScore"
      @submit="submitScore"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { groupApi } from '../api';
import { useGroupStore } from '../stores/group';
import Tournament from './Tournament.vue';
import SubmitScore from './SubmitScore.vue';
import type { Match, Player, PlayerInfo } from '../types';

const props = defineProps<{
  matches: Match[];
  players: Player[];
}>();

const emit = defineEmits<{
  (e: 'create-match', playerIds: string[]): void;
}>();

const groupStore = useGroupStore();
const selectedPlayers = ref([['', ''], ['', '']]);
const selectionError = ref('');
const currentPage = ref(1);
const totalPages = ref(1);
const pageSize = 10;

// Submit Score Modal
const showScoreModal = ref(false);
const selectedMatch = ref<Match | null>(null);

const availablePlayersForSelect = (teamIndex: number, playerIndex: number) => {
  return props.players.filter(player => {
    const selectedIds = selectedPlayers.value.flat().filter((id, idx) => {
      const currentTeamIndex = Math.floor(idx / 2);
      const currentPlayerIndex = idx % 2;
      return id && (teamIndex !== currentTeamIndex || playerIndex !== currentPlayerIndex);
    });
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

const handleCancel = async (match: Match) => {
  if (!confirm('Are you sure you want to cancel this match?')) return;

  try {
    await groupApi.cancelMatch(
      match.group_name,
      match.id,
      groupStore.groupPassword
    );
    await groupStore.loadMatches(currentPage.value, pageSize);
  } catch (error) {
    alert('Failed to cancel match');
  }
};

const handleRandomMatches = async (matches: string[][]) => {
  try {
    await groupApi.createBatchMatches(
      groupStore.currentGroup!.name,
      groupStore.groupPassword,
      matches
    );
    await groupStore.loadMatches(currentPage.value, pageSize);
  } catch (error) {
    alert('Failed to create tournament matches');
  }
};

const changePage = async (page: number) => {
  currentPage.value = page;
  try {
    const response = await groupApi.getMatches(
      groupStore.currentGroup!.name,
      groupStore.groupPassword,
      page,
      pageSize
    );
    totalPages.value = Math.ceil(response.data.total / pageSize);
    await groupStore.loadMatches(page, pageSize);
  } catch (error) {
    alert('Failed to load matches');
  }
};

const showSubmitScore = (match: Match) => {
  selectedMatch.value = match;
  showScoreModal.value = true;
};

const closeSubmitScore = () => {
  showScoreModal.value = false;
  selectedMatch.value = null;
};

const submitScore = async (scoreTeam1: number, scoreTeam2: number) => {
  if (!selectedMatch.value) return;

  try {
    await groupApi.submitResults(
      selectedMatch.value.group_name,
      selectedMatch.value.id,
      groupStore.groupPassword,
      scoreTeam1,
      scoreTeam2
    );
    await Promise.all([
      groupStore.loadMatches(currentPage.value, pageSize),
      groupStore.loadStatistics()
    ]);
  } catch (error) {
    alert('Failed to submit match results');
  }
};
</script>
