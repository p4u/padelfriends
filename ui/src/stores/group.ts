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
    // Save to localStorage
    localStorage.setItem('padel-friends-group', JSON.stringify({ group, password }));
  };

  const clearGroup = () => {
    currentGroup.value = null;
    groupPassword.value = '';
    players.value = [];
    matches.value = [];
    statistics.value = [];
    localStorage.removeItem('padel-friends-group');
  };

  const restoreGroupFromStorage = async (): Promise<boolean> => {
    const storedData = localStorage.getItem('padel-friends-group');
    if (!storedData) return false;

    try {
      const { group, password } = JSON.parse(storedData);
      if (!group || !password) return false;

      // Verify the group still exists and password is valid
      try {
        const response = await groupApi.getByName(group.name, password);
        setGroup(response.data, password);
        return true;
      } catch {
        clearGroup();
        return false;
      }
    } catch {
      clearGroup();
      return false;
    }
  };

  const clearError = () => {
    error.value = null;
  };

  const loadPlayers = async () => {
    if (!currentGroup.value) return;
    loading.value = true;
    error.value = null;
    
    try {
      const response = await groupApi.getPlayers(currentGroup.value.name, groupPassword.value);
      players.value = Array.isArray(response.data) ? markRaw([...response.data]) : [];
    } catch (err) {
      error.value = 'Failed to load players';
      console.error('Failed to load players:', err);
    } finally {
      loading.value = false;
    }
  };

  const loadMatches = async (page: number = 1, pageSize: number = 10) => {
    if (!currentGroup.value) return;
    loading.value = true;
    error.value = null;
    
    try {
      const response = await groupApi.getMatches(currentGroup.value.name, groupPassword.value, page, pageSize);
      matches.value = Array.isArray(response.data.matches) ? markRaw([...response.data.matches]) : [];
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
      const response = await groupApi.getStatistics(currentGroup.value.name, groupPassword.value);
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
    clearGroup,
    restoreGroupFromStorage,
    loadPlayers,
    loadMatches,
    loadStatistics,
    clearError,
  };
});
