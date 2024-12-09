<template>
  <div class="space-y-8">
    <!-- Stats Tabs -->
    <div class="flex justify-center space-x-4">
      <button 
        v-for="tab in tabs" 
        :key="tab.value"
        @click="activeTab = tab.value"
        :class="[
          'px-6 py-3 rounded-lg font-bold transition-all',
          activeTab === tab.value 
            ? 'bg-gradient-to-r from-blue-500 to-indigo-500 text-white shadow-lg' 
            : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
        ]"
      >
        {{ tab.icon }} {{ tab.label }}
      </button>
    </div>

    <!-- Stats Table -->
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="bg-gray-50 dark:bg-gray-700">
              <th class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-white">
                Player
              </th>
              <template v-if="activeTab === 'games'">
                <th 
                  v-for="col in gameColumns" 
                  :key="col.key"
                  @click="sortBy(col.key)"
                  class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-white cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-600"
                >
                  <div class="flex items-center space-x-1">
                    <span>{{ col.label }}</span>
                    <span v-if="sortKey === col.key" class="text-blue-500">
                      {{ sortOrder === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                </th>
              </template>
              <template v-else>
                <th 
                  v-for="col in pointColumns" 
                  :key="col.key"
                  @click="sortBy(col.key)"
                  class="px-6 py-4 text-left text-sm font-semibold text-gray-900 dark:text-white cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-600"
                >
                  <div class="flex items-center space-x-1">
                    <span>{{ col.label }}</span>
                    <span v-if="sortKey === col.key" class="text-blue-500">
                      {{ sortOrder === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                </th>
              </template>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-gray-600">
            <tr 
              v-for="(stat, index) in sortedStats" 
              :key="stat.player_id"
              :class="[
                'hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors',
                index < 3 ? 'font-semibold' : ''
              ]"
            >
              <td class="px-6 py-4">
                <div class="flex items-center space-x-2">
                  <span v-if="index === 0" class="text-2xl">🏆</span>
                  <span v-else-if="index === 1" class="text-2xl">🥈</span>
                  <span v-else-if="index === 2" class="text-2xl">🥉</span>
                  <span :class="index < 3 ? 'text-blue-600 dark:text-blue-400' : 'text-gray-900 dark:text-white'">
                    {{ stat.player_name }}
                  </span>
                </div>
              </td>
              <template v-if="activeTab === 'games'">
                <td 
                  v-for="col in gameColumns" 
                  :key="col.key"
                  class="px-6 py-4 text-gray-900 dark:text-white"
                >
                  <span :class="getValueClass(stat[col.key], col.key)">
                    {{ formatValue(stat[col.key], col.key) }}
                  </span>
                </td>
              </template>
              <template v-else>
                <td 
                  v-for="col in pointColumns" 
                  :key="col.key"
                  class="px-6 py-4 text-gray-900 dark:text-white"
                >
                  <span :class="getValueClass(stat[col.key], col.key)">
                    {{ formatValue(stat[col.key], col.key) }}
                  </span>
                </td>
              </template>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Statistics } from '../types';

const props = defineProps<{
  statistics: Statistics[];
}>();

const activeTab = ref('games');
const sortKey = ref('games_won');
const sortOrder = ref<'asc' | 'desc'>('desc');

const tabs = [
  { value: 'games', label: 'Games', icon: '🎮' },
  { value: 'points', label: 'Points', icon: '📊' }
];

const gameColumns = [
  { key: 'total_games', label: 'Games' },
  { key: 'games_won', label: 'Wins' },
  { key: 'game_win_rate', label: 'Win %' },
  { key: 'games_lost', label: 'Losses' },
  { key: 'game_loss_rate', label: 'Loss %' }
];

const pointColumns = [
  { key: 'total_points', label: 'Points' },
  { key: 'points_won', label: 'Won' },
  { key: 'point_win_rate', label: 'Win %' },
  { key: 'points_lost', label: 'Lost' },
  { key: 'point_loss_rate', label: 'Loss %' }
];

const sortedStats = computed(() => {
  return [...props.statistics].sort((a, b) => {
    const aValue = a[sortKey.value];
    const bValue = b[sortKey.value];
    const modifier = sortOrder.value === 'asc' ? 1 : -1;
    
    if (aValue < bValue) return -1 * modifier;
    if (aValue > bValue) return 1 * modifier;
    return 0;
  });
});

const sortBy = (key: string) => {
  if (sortKey.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortKey.value = key;
    sortOrder.value = 'desc';
  }
};

const formatValue = (value: number, key: string) => {
  if (key.includes('rate')) {
    return `${value.toFixed(1)}%`;
  }
  return value;
};

const getValueClass = (value: number, key: string) => {
  if (!key.includes('rate') && !key.includes('total')) {
    return value > 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400';
  }
  return '';
};
</script>
