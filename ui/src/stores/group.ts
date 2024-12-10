import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { groupApi } from '../api';
import { useSavedGroupsStore } from './savedGroups';
import type { Group, Player, Match, Statistics } from '../types';

export const useGroupStore = defineStore('group', () => {
  const currentGroup = ref<Group | null>(null);
  const groupPassword = ref<string>('');
  const isAuthenticated = ref<boolean>(false);
  const players = ref<Player[]>([]);
  const matches = ref<Match[]>([]);
  const statistics = ref<Statistics | null>(null);
  const savedGroupsStore = useSavedGroupsStore();

  const hasGroup = computed(() => currentGroup.value !== null);

  async function setGroup(group: Group, password?: string) {
    currentGroup.value = group;
    if (password) {
      groupPassword.value = password;
      isAuthenticated.value = true;
      savedGroupsStore.addGroup(group.name, group.name, password);
    } else {
      // Try to restore password from saved groups
      const savedPassword = savedGroupsStore.getGroupPassword(group.name);
      if (savedPassword) {
        groupPassword.value = savedPassword;
        isAuthenticated.value = true;
      } else {
        groupPassword.value = '';
        isAuthenticated.value = false;
      }
    }
  }

  async function authenticate(password: string) {
    if (!currentGroup.value) return false;
    
    try {
      const response = await groupApi.authenticate(currentGroup.value.name, password);
      if (response.data.isAuthenticated) {
        groupPassword.value = password;
        isAuthenticated.value = true;
        savedGroupsStore.addGroup(currentGroup.value.name, currentGroup.value.name, password);
        return true;
      }
    } catch (error) {
      console.error('Authentication failed:', error);
    }
    return false;
  }

  async function loadGroup(name: string) {
    try {
      const response = await groupApi.getByName(name);
      currentGroup.value = response.data;
      
      // Try to restore password from saved groups
      const savedPassword = savedGroupsStore.getGroupPassword(name);
      if (savedPassword) {
        groupPassword.value = savedPassword;
        isAuthenticated.value = true;
      }
      
      return true;
    } catch (error) {
      console.error('Failed to load group:', error);
      return false;
    }
  }

  async function restoreGroupFromStorage() {
    const savedGroups = savedGroupsStore.savedGroups;
    if (savedGroups && savedGroups.length > 0) {
      const mostRecent = savedGroups[0] as { name: string };
      const success = await loadGroup(mostRecent.name);
      return success;
    }
    return false;
  }

  async function clearGroup() {
    if (currentGroup.value?.name) {
      savedGroupsStore.removeGroup(currentGroup.value.name);
    }
    currentGroup.value = null;
    groupPassword.value = '';
    isAuthenticated.value = false;
    players.value = [];
    matches.value = [];
    statistics.value = null;
  }

  async function loadPlayers() {
    if (!currentGroup.value) return;
    
    try {
      const response = await groupApi.getPlayers(currentGroup.value.name);
      players.value = response.data;
    } catch (error) {
      console.error('Failed to load players:', error);
    }
  }

  async function loadMatches(page: number = 1, pageSize: number = 10) {
    if (!currentGroup.value) return;
    
    try {
      const response = await groupApi.getMatches(currentGroup.value.name, page, pageSize);
      matches.value = response.data.matches;
    } catch (error) {
      console.error('Failed to load matches:', error);
    }
  }

  async function loadStatistics() {
    if (!currentGroup.value) return;
    
    try {
      const response = await groupApi.getStatistics(currentGroup.value.name);
      statistics.value = response.data;
    } catch (error) {
      console.error('Failed to load statistics:', error);
    }
  }

  return {
    currentGroup,
    groupPassword,
    isAuthenticated,
    players,
    matches,
    statistics,
    hasGroup,
    setGroup,
    authenticate,
    loadGroup,
    restoreGroupFromStorage,
    clearGroup,
    loadPlayers,
    loadMatches,
    loadStatistics,
  };
});
