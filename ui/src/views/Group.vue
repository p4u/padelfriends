<template>
  <div class="max-w-4xl mx-auto space-y-8">
    <!-- Group Header -->
    <div class="modern-container bg-white dark:bg-gray-800">
      <h1 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
        🎾 {{ currentGroup?.name }} 🏸
      </h1>
      <div class="flex flex-wrap gap-2">
        <button 
          v-for="tab in tabs" 
          :key="tab.value"
          @click="activeTab = tab.value" 
          :class="[
            'modern-button',
            activeTab === tab.value ? 'bg-primary ring-2 ring-primary ring-offset-2' : 'bg-gray-500'
          ]"
        >
          {{ tab.icon }} {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- Content Tabs -->
    <div class="modern-container bg-white dark:bg-gray-800">
      <div v-if="loading" class="text-center text-gray-900 dark:text-white">
        Loading... 🔄
      </div>
      <div v-else-if="error" class="text-center text-red-500">
        {{ error }} 😢
      </div>
      <component
        v-else
        :is="currentComponent"
        v-bind="componentProps"
        @add-player="addPlayer"
        @create-match="createMatch"
        @submit-score="showSubmitScore"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useGroupStore } from '../stores/group';
import { groupApi } from '../api';
import { PlayerList, MatchList, StatisticsList } from '../components';
import type { Match } from '../types';

const route = useRoute();
const router = useRouter();
const groupStore = useGroupStore();
const activeTab = ref('players');
const loading = ref(false);
const error = ref<string | null>(null);

const tabs = [
  { value: 'players', label: 'Players', icon: '🏃' },
  { value: 'matches', label: 'Matches', icon: '🎾' },
  { value: 'statistics', label: 'Stats', icon: '📊' }
];

const currentGroup = computed(() => groupStore.currentGroup);
const players = computed(() => groupStore.players);
const matches = computed(() => groupStore.matches);
const statistics = computed(() => groupStore.statistics);

const currentComponent = computed(() => {
  switch (activeTab.value) {
    case 'players': return PlayerList;
    case 'matches': return MatchList;
    case 'statistics': return StatisticsList;
    default: return PlayerList;
  }
});

const componentProps = computed(() => {
  switch (activeTab.value) {
    case 'players': return { players: players.value };
    case 'matches': return { matches: matches.value, players: players.value };
    case 'statistics': return { statistics: statistics.value };
    default: return {};
  }
});

onMounted(async () => {
  if (!currentGroup.value) {
    const restored = await groupStore.restoreGroupFromStorage();
    if (!restored) {
      router.push('/');
      return;
    }
  }

  loading.value = true;
  error.value = null;

  try {
    await Promise.all([
      groupStore.loadPlayers(),
      groupStore.loadMatches(),
      groupStore.loadStatistics()
    ]);
  } catch (err) {
    error.value = 'Failed to load group data';
    console.error('Failed to load group data:', err);
  } finally {
    loading.value = false;
  }
});

const addPlayer = async (name: string) => {
  if (!currentGroup.value) return;
  
  try {
    await groupApi.addPlayer(
      currentGroup.value.name,
      groupStore.groupPassword,
      name
    );
    await groupStore.loadPlayers();
  } catch (error) {
    alert('❌ Failed to add player');
  }
};

const createMatch = async (playerIds: string[]) => {
  if (!currentGroup.value) return;
  
  try {
    await groupApi.createMatch(
      currentGroup.value.name,
      groupStore.groupPassword,
      playerIds
    );
    await groupStore.loadMatches();
  } catch (error) {
    alert('❌ Failed to create match');
  }
};

const showSubmitScore = async (match: Match) => {
  const score1 = prompt('🎾 Enter score for Team 1:');
  const score2 = prompt('🎾 Enter score for Team 2:');
  
  if (score1 === null || score2 === null) return;
  
  const scoreTeam1 = parseInt(score1);
  const scoreTeam2 = parseInt(score2);
  
  if (isNaN(scoreTeam1) || isNaN(scoreTeam2)) {
    alert('❌ Please enter valid scores');
    return;
  }
  
  try {
    await groupApi.submitResults(
      match.group_name,
      match.id,
      groupStore.groupPassword,
      scoreTeam1,
      scoreTeam2
    );
    await Promise.all([
      groupStore.loadMatches(),
      groupStore.loadStatistics()
    ]);
  } catch (error) {
    alert('❌ Failed to submit match results');
  }
};
</script>
