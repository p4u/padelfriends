import { defineStore } from 'pinia';
import { ref, markRaw } from 'vue';
import type { Group, Player, Match, Statistics } from '../types';
import { groupApi } from '../api';

export const useGroupStore = defineStore('group', () => {
  const currentGroup = ref<Group | null>(null);
  const players = ref<Player[]>([]);
  const matches = ref<Match[]>([]);
  const statistics = ref<Statistics[]>([]);
  const groupPassword = ref<string>('');
  const loading = ref(false);
  const error = ref<string | null>(null);

  const setGroup = (group: Group, password: string) => {
    currentGroup.value = markRaw({ ...group });
    groupPassword.value = password;
  };

  const clearError = () => {
    error.value = null;
  };

  const loadPlayers = async () => {
    if (!currentGroup.value) return;
    loading.value = true;
    error.value = null;
    
    try {
      const response = await groupApi.getPlayers(currentGroup.value.id, groupPassword.value);
      players.value = markRaw([...response.data]);
    } catch (err) {
      error.value = 'Failed to load players';
      console.error('Failed to load players:', err);
    } finally {
      loading.value = false;
    }
  };

  const loadMatches = async () => {
    if (!currentGroup.value) return;
    loading.value = true;
    error.value = null;
    
    try {
      const response = await groupApi.getMatches(currentGroup.value.id, groupPassword.value);
      matches.value = markRaw([...response.data]);
    } catch (err) {
      error.value = 'Failed to load matches';
      console.error('Failed to load matches:', err);
    } finally {
      loading.value = false;
    }
  };

  const loadStatistics = async () => {
    if (!currentGroup.value) return;
    loading.value = true;
    error.value = null;
    
    try {
      const response = await groupApi.getStatistics(currentGroup.value.id, groupPassword.value);
      statistics.value = markRaw([...response.data]);
    } catch (err) {
      error.value = 'Failed to load statistics';
      console.error('Failed to load statistics:', err);
    } finally {
      loading.value = false;
    }
  };

  return {
    currentGroup,
    players,
    matches,
    statistics,
    groupPassword,
    loading,
    error,
    setGroup,
    loadPlayers,
    loadMatches,
    loadStatistics,
    clearError,
  };
});