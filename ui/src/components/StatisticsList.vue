<template>
  <div class="space-y-4">
    <div v-for="stat in statistics" :key="stat.player_id" class="p-4 bg-gray-100 rounded">
      <h3 class="font-bold">{{ stat.player_name }}</h3>
      <div class="grid grid-cols-3 gap-4 mt-2">
        <div>
          <p class="text-sm">Matches Played</p>
          <p class="font-bold">{{ stat.matches_played }}</p>
        </div>
        <div>
          <p class="text-sm">Matches Won</p>
          <p class="font-bold">{{ stat.matches_won }}</p>
        </div>
        <div>
          <p class="text-sm">Win Rate</p>
          <p class="font-bold">{{ (stat.win_rate * 100).toFixed(1) }}%</p>
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