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
    // Remove any existing group with same id or name
    const filtered = savedGroups.value.filter(g => g.id !== id && g.name !== name);

    // Add new group at the start
    savedGroups.value = [
      {
        id,
        name,
        password,
        lastAccessed: new Date().toISOString()
      },
      ...filtered
    ];
  };

  const removeGroup = (id: string) => {
    savedGroups.value = savedGroups.value.filter(g => g.id !== id);
  };

  const getGroupPassword = (id: string) => {
    const group = savedGroups.value.find(g => g.id === id);
    return group?.password || null;
  };

  return {
    savedGroups,
    addGroup,
    removeGroup,
    getGroupPassword
  };
});
