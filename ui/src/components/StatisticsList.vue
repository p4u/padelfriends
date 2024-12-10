<template>
  <div class="space-y-6">
    <!-- Stats Tabs -->
    <div class="flex justify-center space-x-4">
      <button 
        v-for="tab in tabs" 
        :key="tab.value"
        @click="activeTab = tab.value"
        :class="[
          'px-4 py-2 rounded-lg font-bold transition-all',
          activeTab === tab.value 
            ? 'bg-gradient-to-r from-blue-500 to-indigo-500 text-white shadow-lg' 
            : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
        ]"
      >
        {{ tab.icon }} {{ tab.label }}
      </button>
    </div>

    <!-- Stats Table -->
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden min-h-[400px]">
      <div v-if="!statistics?.length" class="flex items-center justify-center h-[400px] text-gray-500 dark:text-gray-400">
        No statistics available
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full table-fixed">
          <thead>
            <tr class="bg-gray-50 dark:bg-gray-700">
              <th class="w-1/3 px-2 py-2 text-left text-sm font-semibold text-gray-900 dark:text-white">
                Player
              </th>
              <template v-if="activeTab === 'sets'">
                <th 
                  v-for="col in setColumns" 
                  :key="col.key"
                  @click="sortBy(col.key)"
                  class="w-[22%] px-2 py-2 text-center text-sm font-semibold text-gray-900 dark:text-white cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-600"
                >
                  {{ col.label }}
                  <span v-if="sortKey === col.key" class="text-blue-500 ml-1">
                    {{ sortOrder === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
              </template>
              <template v-else>
                <th 
                  v-for="col in pointColumns" 
                  :key="col.key"
                  @click="sortBy(col.key)"
                  class="w-[22%] px-2 py-2 text-center text-sm font-semibold text-gray-900 dark:text-white cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-600"
                >
                  {{ col.label }}
                  <span v-if="sortKey === col.key" class="text-blue-500 ml-1">
                    {{ sortOrder === 'asc' ? '↑' : '↓' }}
                  </span>
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
              <td class="px-2 py-2 truncate">
                <div class="flex items-center gap-1">
                  <span v-if="index === 0" class="text-lg">🏆</span>
                  <span v-else-if="index === 1" class="text-lg">🥈</span>
                  <span v-else-if="index === 2" class="text-lg">🥉</span>
                  <span :class="index < 3 ? 'text-blue-600 dark:text-blue-400' : 'text-gray-900 dark:text-white'">
                    {{ stat.player_name }}
                  </span>
                </div>
              </td>
              <template v-if="activeTab === 'sets'">
                <td 
                  v-for="col in setColumns" 
                  :key="col.key"
                  class="px-2 py-2 text-center text-sm"
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
                  class="px-2 py-2 text-center text-sm"
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

const activeTab = ref('sets');
const sortKey = ref('games_won');
const sortOrder = ref<'asc' | 'desc'>('desc');

const tabs = [
  { value: 'sets', label: 'Sets', icon: '🎮' },
  { value: 'points', label: 'Points', icon: '📊' }
];

const setColumns = [
  { key: 'games_won', label: 'W' },
  { key: 'games_lost', label: 'L' },
  { key: 'game_win_rate', label: 'W%' }
];

const pointColumns = [
  { key: 'points_won', label: 'W' },
  { key: 'points_lost', label: 'L' },
  { key: 'point_win_rate', label: 'W%' }
];

const sortedStats = computed(() => {
  if (!props.statistics?.length) return [];
  
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
    return `${value.toFixed(0)}%`;
  }
  return value;
};

const getValueClass = (value: number, key: string) => {
  if (key.includes('lost')) {
    return 'text-red-600 dark:text-red-400 font-semibold';
  }
  if (!key.includes('rate') && !key.includes('total') && !key.includes('lost')) {
    return value > 0 ? 'text-green-600 dark:text-green-400 font-semibold' : '';
  }
  return '';
};
</script>
