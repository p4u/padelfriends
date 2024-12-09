<template>
  <div class="space-y-8">
    <!-- Create Match Form -->
    <form @submit.prevent="handleSubmit" class="space-y-4 bg-gradient-to-r from-blue-500 to-indigo-500 p-6 rounded-xl shadow-lg">
      <h3 class="text-xl font-bold text-white mb-4">{{ t('matches.newMatch') }}</h3>
      <div class="grid grid-cols-2 gap-6">
        <div v-for="(team, index) in [1, 2]" :key="team" class="space-y-3">
          <div class="flex items-center">
            <span class="text-white font-medium">{{ t(`matches.team${team}`) }}</span>
          </div>
          <div class="space-y-2">
            <select
              :value="getPlayerValue(index-1, 0)"
              @input="updatePlayer(index-1, 0, $event)"
              class="modern-input w-full text-gray-900 dark:text-white bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
              required
            >
              <option value="">{{ t('matches.selectPlayer') }} 1</option>
              <option 
                v-for="player in getAvailablePlayers(index-1, 0)" 
                :key="player.id" 
                :value="player.id"
              >
                {{ player.name }}
              </option>
            </select>
            <select
              :value="getPlayerValue(index-1, 1)"
              @input="updatePlayer(index-1, 1, $event)"
              class="modern-input w-full text-gray-900 dark:text-white bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
              required
            >
              <option value="">{{ t('matches.selectPlayer') }} 2</option>
              <option 
                v-for="player in getAvailablePlayers(index-1, 1)" 
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
        {{ t('matches.startSingleMatch') }}
      </button>
    </form>

    <!-- Tournament Generator -->
    <Tournament 
      :players="players"
      @create-matches="handleRandomMatches"
    />

    <!-- Matches List -->
    <div class="mt-12 border-t border-gray-200 dark:border-gray-700 pt-8">
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
            {{ t('matches.matchHistory') }}
          </h2>
        </div>
        
        <div class="p-6 space-y-4">
          <div v-for="match in matches" :key="match.id" 
               class="bg-gray-50 dark:bg-gray-700 rounded-xl p-6 hover:shadow-md transition-shadow">
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
                  <div class="px-4 py-2 rounded-full bg-gray-100 dark:bg-gray-600 font-bold text-gray-600 dark:text-gray-300">
                    {{ t('matches.vs') }}
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
                      @click="$emit('submit-score', match)"
                      class="modern-button bg-gradient-to-r from-green-500 to-emerald-500 text-white font-bold"
                    >
                      {{ t('matches.submitScore') }}
                    </button>
                    <button
                      @click="handleCancel(match)"
                      class="modern-button bg-gradient-to-r from-red-500 to-pink-500 text-white font-bold"
                    >
                      {{ t('matches.cancelMatch') }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="px-6 py-4 border-t border-gray-200 dark:border-gray-700">
          <div class="flex justify-center space-x-2">
            <button 
              @click="changePage(currentPage - 1)"
              :disabled="currentPage === 1"
              class="px-4 py-2 rounded-lg bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 disabled:opacity-50"
            >
              {{ t('common.previous') }}
            </button>
            <div class="flex items-center space-x-2">
              <span class="text-gray-600 dark:text-gray-400">
                {{ t('common.page', { current: currentPage, total: totalPages }) }}
              </span>
            </div>
            <button 
              @click="changePage(currentPage + 1)"
              :disabled="currentPage === totalPages"
              class="px-4 py-2 rounded-lg bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 disabled:opacity-50"
            >
              {{ t('common.next') }}
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
import { groupApi } from '../api';
import { useGroupStore } from '../stores/group';
import Tournament from './Tournament.vue';
import type { Match, Player, PlayerInfo } from '../types';

const props = defineProps<{
  matches: Match[];
  players: Player[];
}>();

const emit = defineEmits<{
  (e: 'create-match', playerIds: string[]): void;
  (e: 'submit-score', match: Match): void;
}>();

const { t } = useI18n();
const groupStore = useGroupStore();
const playerSelections = ref<string[][]>([
  ['', ''],
  ['', '']
]);
const selectionError = ref('');
const currentPage = ref(1);
const totalPages = ref(1);
const pageSize = 10;

const getPlayerValue = (teamIndex: number, playerIndex: number): string => {
  return playerSelections.value?.[teamIndex]?.[playerIndex] || '';
};

const updatePlayer = (teamIndex: number, playerIndex: number, event: Event) => {
  const value = (event.target as HTMLSelectElement).value;
  if (!playerSelections.value[teamIndex]) {
    playerSelections.value[teamIndex] = ['', ''];
  }
  playerSelections.value[teamIndex][playerIndex] = value;
  validateSelection();
};

const getAvailablePlayers = (teamIndex: number, playerIndex: number) => {
  const currentValue = getPlayerValue(teamIndex, playerIndex);
  const selectedIds = playerSelections.value
    .flatMap(team => team || ['', ''])
    .filter(id => id !== '' && id !== currentValue);

  return props.players.filter(player => 
    !selectedIds.includes(player.id) || player.id === currentValue
  );
};

const validateSelection = () => {
  const selectedIds = playerSelections.value
    .flatMap(team => team || ['', ''])
    .filter(id => id !== '');
  
  const uniqueIds = new Set(selectedIds);
  
  if (selectedIds.length !== uniqueIds.size) {
    selectionError.value = t('errors.duplicatePlayer');
  } else {
    selectionError.value = '';
  }
};

const getTeamNames = (team: PlayerInfo[]) => {
  if (!team || !Array.isArray(team)) return '';
  return team.map(player => player?.name || '').filter(Boolean).join(' & ');
};

const formatDate = (timestamp: string) => {
  const date = new Date(timestamp);
  const now = new Date();
  const diffTime = Math.abs(now.getTime() - date.getTime());
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

  if (diffDays === 0) {
    return t('common.today') + ', ' + date.toLocaleTimeString();
  } else if (diffDays === 1) {
    return t('common.yesterday') + ', ' + date.toLocaleTimeString();
  } else {
    return date.toLocaleDateString();
  }
};

const handleSubmit = () => {
  if (selectionError.value) return;
  
  const playerIds = playerSelections.value.flatMap(team => team || ['', '']).filter(id => id !== '');
  if (playerIds.length === 4) {
    emit('create-match', playerIds);
    playerSelections.value = [['', ''], ['', '']];
    selectionError.value = '';
  }
};

const handleCancel = async (match: Match) => {
  if (!confirm(t('matches.confirmCancel'))) return;

  try {
    await groupApi.cancelMatch(
      match.group_name,
      match.id,
      groupStore.groupPassword
    );
    await groupStore.loadMatches(currentPage.value, pageSize);
  } catch (error) {
    alert(t('errors.failedToDelete'));
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
    alert(t('errors.failedToCreate'));
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
    alert(t('errors.failedToLoad'));
  }
};
</script>
