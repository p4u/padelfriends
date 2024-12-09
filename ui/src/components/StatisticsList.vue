<template>
  <div class="space-y-8">
    <!-- Statistics Header -->
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
          {{ t('statistics.title') }}
        </h2>
      </div>

      <div v-if="!statistics.length" class="p-6 text-center text-gray-500 dark:text-gray-400">
        {{ t('statistics.noStats') }}
      </div>
      
      <div v-else class="p-6">
        <!-- Games Statistics -->
        <div class="mb-8">
          <h3 class="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
            {{ t('statistics.games') }}
          </h3>
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead>
                <tr>
                  <th 
                    v-for="col in gameColumns" 
                    :key="col.key"
                    @click="sortBy(col.key)"
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700"
                  >
                    {{ t(`statistics.${col.key}`) }}
                    <span v-if="sortKey === col.key" class="ml-1">
                      {{ sortOrder === 'asc' ? '↑' : '↓' }}
                    </span>
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                <tr v-for="stat in sortedStats" :key="stat.player_id">
                  <td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-white">
                    {{ stat.player_name }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-white">
                    {{ stat.total_games }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span class="text-green-600 dark:text-green-400">
                      {{ stat.games_won }} ({{ formatPercent(stat.game_win_rate) }})
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span class="text-red-600 dark:text-red-400">
                      {{ stat.games_lost }} ({{ formatPercent(stat.game_loss_rate) }})
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Points Statistics -->
        <div>
          <h3 class="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
            {{ t('statistics.points') }}
          </h3>
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead>
                <tr>
                  <th 
                    v-for="col in pointColumns" 
                    :key="col.key"
                    @click="sortBy(col.key)"
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700"
                  >
                    {{ t(`statistics.${col.key}`) }}
                    <span v-if="sortKey === col.key" class="ml-1">
                      {{ sortOrder === 'asc' ? '↑' : '↓' }}
                    </span>
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                <tr v-for="stat in sortedStats" :key="stat.player_id">
                  <td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-white">
                    {{ stat.player_name }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-gray-900 dark:text-white">
                    {{ stat.total_points }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span class="text-green-600 dark:text-green-400">
                      {{ stat.points_won }} ({{ formatPercent(stat.point_win_rate) }})
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span class="text-red-600 dark:text-red-400">
                      {{ stat.points_lost }} ({{ formatPercent(stat.point_loss_rate) }})
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from '../i18n';
import type { Statistics } from '../types';

const props = defineProps<{
  statistics: Statistics[];
}>();

const { t } = useI18n();
const sortKey = ref('games_won');
const sortOrder = ref<'asc' | 'desc'>('desc');

const gameColumns = [
  { key: 'player_name', label: 'Player' },
  { key: 'total_games', label: 'Total Games' },
  { key: 'games_won', label: 'Games Won' },
  { key: 'games_lost', label: 'Games Lost' }
];

const pointColumns = [
  { key: 'player_name', label: 'Player' },
  { key: 'total_points', label: 'Total Points' },
  { key: 'points_won', label: 'Points Won' },
  { key: 'points_lost', label: 'Points Lost' }
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

const formatPercent = (value: number) => {
  return `${value.toFixed(1)}%`;
};
</script>
