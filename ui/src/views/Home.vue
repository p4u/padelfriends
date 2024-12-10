<template>
  <div class="max-w-4xl mx-auto space-y-8 p-4">
    <div class="space-y-8">
      <!-- Header -->
      <div class="text-center space-y-4">
        <h1 class="text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-500 to-indigo-500">
          Padel Friends
        </h1>
        <p class="text-gray-600 dark:text-gray-300">Connect, Play, Enjoy</p>
        <div class="flex justify-center">
          <img src="/icons/icon-512x512.png" alt="Padel Friends Logo" class="w-32 h-32" />
        </div>
      </div>
      
      <!-- Recent Groups -->
      <div v-if="savedGroups.length > 0" class="modern-container min-h-[200px]">
        <h2 class="text-2xl font-semibold text-gray-900 dark:text-white text-center mb-6">Recent Groups</h2>
        <div class="grid gap-4 sm:grid-cols-2">
          <div 
            v-for="group in savedGroups" 
            :key="group.id"
            class="bg-gray-50 dark:bg-gray-700 rounded-xl shadow-md hover:shadow-lg transition-all cursor-pointer p-4"
            @click="joinSavedGroup(group)"
          >
            <div class="flex justify-between items-center">
              <div class="flex items-center space-x-3">
                <span class="text-2xl">🎾</span>
                <span class="font-semibold text-gray-900 dark:text-white transition-colors duration-200 hover:text-blue-600 dark:hover:text-blue-400">
                  {{ group.name }}
                </span>
              </div>
              <button 
                class="text-red-500 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300 p-2"
                @click.stop="removeGroup(group)"
              >
                ❌
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Available Groups -->
      <div v-if="filteredGroups.length > 0" class="modern-container min-h-[200px]">
        <h2 class="text-2xl font-semibold text-gray-900 dark:text-white text-center mb-6">Available Groups</h2>
        <div class="grid gap-4 sm:grid-cols-2">
          <div 
            v-for="group in filteredGroups" 
            :key="group.id"
            class="bg-gray-50 dark:bg-gray-700 rounded-xl shadow-md hover:shadow-lg transition-all cursor-pointer p-4"
            @click="joinGroup(group)"
          >
            <div class="flex items-center space-x-3">
              <span class="text-2xl">🎯</span>
              <span class="font-semibold text-gray-900 dark:text-white transition-colors duration-200 hover:text-blue-600 dark:hover:text-blue-400">
                {{ group.name }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Register Group Form -->
      <div class="modern-container min-h-[200px] bg-gradient-to-r from-blue-500 to-indigo-500">
        <div class="text-center mb-6">
          <h2 class="text-2xl font-bold text-white">Register a new group</h2>
          <p class="text-white/80 mt-2">for casual play, a tournament or a league, it's free!</p>
        </div>
        <form @submit.prevent="createGroup" class="space-y-4">
          <input
            v-model="form.name"
            type="text"
            placeholder="Enter a new group name"
            class="modern-input w-full bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
            required
          />
          <input
            v-model="form.password"
            type="password"
            placeholder="Enter a password for you and your friends"
            class="modern-input w-full bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
            required
          />
          <button type="submit" class="modern-button w-full bg-gradient-to-r from-green-500 to-emerald-500 text-white font-bold hover:from-green-600 hover:to-emerald-600 transition-all duration-200">
            🎮 Create Group
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import { groupApi } from '../api';
import { useGroupStore } from '../stores/group';
import { useSavedGroupsStore } from '../stores/savedGroups';
import type { Group } from '../types';

const router = useRouter();
const groupStore = useGroupStore();
const savedGroupsStore = useSavedGroupsStore();

const groups = ref<Group[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const form = ref({ name: '', password: '' });

const { savedGroups } = storeToRefs(savedGroupsStore);

const filteredGroups = computed(() => {
  const savedGroupNames = new Set(savedGroups.value.map(group => group.name));
  return groups.value.filter(group => !savedGroupNames.has(group.name));
});

onMounted(async () => {
  try {
    const response = await groupApi.list();
    if (Array.isArray(response.data)) {
      groups.value = response.data.map(group => ({
        id: group.id || group.name, // Fallback to name if id is not available
        name: group.name,
        created_at: group.created_at
      }));
    } else {
      error.value = 'Invalid data format received from server';
    }
  } catch (err) {
    error.value = typeof err === 'string' ? err : 'Failed to fetch groups';
  } finally {
    loading.value = false;
  }
});

const joinGroup = async (group: Group) => {
  if (!group || !group.name) {
    console.error('Invalid group data:', group);
    return;
  }

  try {
    const response = await groupApi.getByName(group.name);
    const joinedGroup = response.data;
    await groupStore.setGroup(joinedGroup);
    router.push(`/group/${joinedGroup.name}`);
  } catch (err) {
    alert('Failed to join group');
  }
};

const createGroup = async () => {
  try {
    const response = await groupApi.create(form.value.name, form.value.password);
    const group = response.data;
    await groupStore.setGroup(group, form.value.password);
    // Use group.id if available, otherwise use name
    savedGroupsStore.addGroup(group.id || group.name, group.name, form.value.password);
    router.push(`/group/${group.name}`);
  } catch (err) {
    alert('Failed to create group. Please try again.');
  }
};

const joinSavedGroup = async (group: { id: string, name: string, password: string }) => {
  try {
    const response = await groupApi.getByName(group.name);
    const joinedGroup = response.data;
    await groupStore.setGroup(joinedGroup, group.password);
    router.push(`/group/${group.name}`);
  } catch (err) {
    alert('Failed to join group. The group may no longer exist.');
    savedGroupsStore.removeGroup(group.id);
  }
};

const removeGroup = (group: { id: string, name: string }) => {
  if (confirm('Are you sure you want to remove this group from your recent list?')) {
    savedGroupsStore.removeGroup(group.id);
  }
};
</script>
