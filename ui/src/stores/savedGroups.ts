import { defineStore } from 'pinia';
import { useLocalStorage } from '@vueuse/core';

interface SavedGroup {
  id: string;
  name: string;
  password: string;
  lastAccessed: string;
}

export const useSavedGroupsStore = defineStore('savedGroups', () => {
  const savedGroups = useLocalStorage<SavedGroup[]>('padel-friends-groups', []);

  const addGroup = (id: string, name: string, password: string) => {
    const group = {
      id,
      name,
      password,
      lastAccessed: new Date().toISOString()
    };

    const newGroups = savedGroups.value.filter(g => g.id !== id);
    savedGroups.value = [...newGroups, group];
  };

  const removeGroup = (id: string) => {
    savedGroups.value = savedGroups.value.filter(g => g.id !== id);
  };

  const getGroupPassword = (id: string) => {
    const group = savedGroups.value.find(g => g.id === id);
    return group?.password;
  };

  return {
    savedGroups,
    addGroup,
    removeGroup,
    getGroupPassword,
  };
});