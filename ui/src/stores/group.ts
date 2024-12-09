import { defineStore } from 'pinia';
import { ref } from 'vue';
import { groupApi } from '../api';
import type { Group, Player, Match, Statistics } from '../types';

interface SavedGroup {
  name: string;
  password: string;
}

export const useGroupStore = defineStore('group', () => {
  const currentGroup = ref<Group | null>(null);
  const groupPassword = ref<string>('');
  const players = ref<Player[]>([]);
  const matches = ref<Match[]>([]);
  const statistics = ref<Statistics[]>([]);

  const setGroup = async (name: string, password: string) => {
    try {
      const response = await groupApi.getGroup(name, password);
      currentGroup.value = response.data;
      groupPassword.value = password;

      // Save to localStorage
      const savedGroups = getSavedGroups();
      const existingIndex = savedGroups.findIndex((g: SavedGroup) => g.name === name);
      
      if (existingIndex !== -1) {
        savedGroups[existingIndex] = { name, password };
      } else {
        savedGroups.push({ name, password });
      }
      
      localStorage.setItem('savedGroups', JSON.stringify(savedGroups));
      
      return true;
    } catch (error) {
      console.error('Error setting group:', error);
      return false;
    }
  };

  const getSavedGroups = (): SavedGroup[] => {
    try {
      const savedGroupsStr = localStorage.getItem('savedGroups');
      return savedGroupsStr ? JSON.parse(savedGroupsStr) : [];
    } catch (error) {
      console.error('Error parsing saved groups:', error);
      return [];
    }
  };

  const restoreGroupFromStorage = async () => {
    const savedGroups = getSavedGroups();
    const currentPath = window.location.pathname;
    const groupName = currentPath.split('/').pop();

    if (groupName) {
      const savedGroup = savedGroups.find((g: SavedGroup) => g.name === groupName);
      if (savedGroup) {
        return await setGroup(savedGroup.name, savedGroup.password);
      }
    }
    return false;
  };

  const loadPlayers = async () => {
    if (!currentGroup.value || !groupPassword.value) return;
    
    try {
      const response = await groupApi.getPlayers(
        currentGroup.value.name,
        groupPassword.value
      );
      players.value = response.data;
    } catch (error) {
      console.error('Error loading players:', error);
      throw error;
    }
  };

  const loadMatches = async (page: number = 1, pageSize: number = 10) => {
    if (!currentGroup.value || !groupPassword.value) return;
    
    try {
      const response = await groupApi.getMatches(
        currentGroup.value.name,
        groupPassword.value,
        page,
        pageSize
      );
      matches.value = response.data.matches;
    } catch (error) {
      console.error('Error loading matches:', error);
      throw error;
    }
  };

  const loadStatistics = async () => {
    if (!currentGroup.value || !groupPassword.value) return;
    
    try {
      const response = await groupApi.getStatistics(
        currentGroup.value.name,
        groupPassword.value
      );
      statistics.value = response.data;
    } catch (error) {
      console.error('Error loading statistics:', error);
      throw error;
    }
  };

  return {
    currentGroup,
    groupPassword,
    players,
    matches,
    statistics,
    setGroup,
    restoreGroupFromStorage,
    loadPlayers,
    loadMatches,
    loadStatistics
  };
});
