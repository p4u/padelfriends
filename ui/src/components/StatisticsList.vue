<template>
  <div class="space-y-4">
    <div v-for="stat in statistics" :key="stat.player_id" class="p-4 bg-white dark:bg-gray-800 rounded-lg shadow-sm">
      <h3 class="font-semibold text-gray-900 dark:text-white">{{ stat.player_name }}</h3>
      <div class="grid grid-cols-3 gap-4 mt-2">
        <div>
          <p class="text-sm text-gray-600 dark:text-gray-300">Matches Played</p>
          <p class="font-bold text-gray-900 dark:text-white">{{ stat.matches_played }}</p>
        </div>
        <div>
          <p class="text-sm text-gray-600 dark:text-gray-300">Matches Won</p>
          <p class="font-bold text-gray-900 dark:text-white">{{ stat.matches_won }}</p>
        </div>
        <div>
          <p class="text-sm text-gray-600 dark:text-gray-300">Win Rate</p>
          <p class="font-bold text-gray-900 dark:text-white">{{ (stat.win_rate * 100).toFixed(1) }}%</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Statistics } from '../types';

defineProps<{
  statistics: Statistics[];
}>();
</script>
