<template>
  <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
    <div class="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-md w-full p-6">
      <div class="space-y-6">
        <!-- Header -->
        <div class="flex justify-between items-center">
          <h3 class="text-xl font-bold text-gray-900 dark:text-white">
            {{ t('matches.submitScore') }}
          </h3>
          <button @click="close" class="text-gray-500 hover:text-gray-700 dark:text-gray-400">
            ❌
          </button>
        </div>

        <!-- Teams and Scores -->
        <div class="space-y-6">
          <!-- Team 1 -->
          <div class="space-y-2">
            <label class="block font-medium text-gray-900 dark:text-white">
              {{ getTeamNames(match.team1) }}
            </label>
            <div class="flex items-center space-x-4">
              <input 
                type="range" 
                v-model="scoreTeam1" 
                min="0" 
                max="10" 
                class="flex-1"
              >
              <span class="text-2xl font-bold text-blue-600 dark:text-blue-400 w-12 text-center">
                {{ scoreTeam1 }}
              </span>
            </div>
          </div>

          <!-- VS Divider -->
          <div class="flex items-center justify-center">
            <span class="px-4 py-2 rounded-full bg-gray-100 dark:bg-gray-700 font-bold text-gray-600 dark:text-gray-300">
              {{ t('matches.vs') }}
            </span>
          </div>

          <!-- Team 2 -->
          <div class="space-y-2">
            <label class="block font-medium text-gray-900 dark:text-white">
              {{ getTeamNames(match.team2) }}
            </label>
            <div class="flex items-center space-x-4">
              <input 
                type="range" 
                v-model="scoreTeam2" 
                min="0" 
                max="10" 
                class="flex-1"
              >
              <span class="text-2xl font-bold text-indigo-600 dark:text-indigo-400 w-12 text-center">
                {{ scoreTeam2 }}
              </span>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex justify-end space-x-3 pt-4">
          <button 
            @click="close"
            class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 dark:text-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600"
          >
            {{ t('common.cancel') }}
          </button>
          <button 
            @click="submit"
            class="px-4 py-2 text-white bg-gradient-to-r from-green-500 to-emerald-500 rounded-lg hover:shadow-lg"
          >
            {{ t('common.submit') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from '../i18n';
import type { Match, PlayerInfo } from '../types';

const props = defineProps<{
  show: boolean;
  match: Match;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'submit', scoreTeam1: number, scoreTeam2: number): void;
}>();

const { t } = useI18n();
const scoreTeam1 = ref(0);
const scoreTeam2 = ref(0);

const getTeamNames = (team: PlayerInfo[]) => {
  if (!team || !Array.isArray(team)) return '';
  return team.map(player => player?.name || '').filter(Boolean).join(' & ');
};

const close = () => {
  scoreTeam1.value = 0;
  scoreTeam2.value = 0;
  emit('close');
};

const submit = () => {
  emit('submit', scoreTeam1.value, scoreTeam2.value);
  close();
};
</script>
