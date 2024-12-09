<template>
  <div class="max-w-4xl mx-auto space-y-8 p-4">
    <div class="space-y-8">
      <!-- Header -->
      <div class="text-center space-y-4">
        <h1 class="text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-500 to-indigo-500">
          Padel Friends
        </h1>
        <p class="text-gray-600 dark:text-gray-300">Connect, Play, Compete</p>
      </div>
      
      <!-- Recent Groups -->
      <div v-if="savedGroups.length > 0" class="space-y-4">
        <h2 class="text-2xl font-semibold text-gray-900 dark:text-white">Recent Groups</h2>
        <div class="grid gap-4 sm:grid-cols-2">
          <div 
            v-for="group in savedGroups" 
            :key="group.id"
            class="bg-white dark:bg-gray-800 rounded-xl shadow-md hover:shadow-lg transition-all cursor-pointer p-4"
            @click="joinSavedGroup(group)"
          >
            <div class="flex justify-between items-center">
              <div class="flex items-center space-x-3">
                <span class="text-2xl">🎾</span>
                <span class="font-semibold text-gray-900 dark:text-white">{{ group.name }}</span>
              </div>
              <button 
                class="text-red-500 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300 p-2"
                @click.stop="removeGroup(group.id)"
              >
                ❌
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Available Groups -->
      <div v-if="filteredGroups.length > 0" class="space-y-4">
        <h2 class="text-2xl font-semibold text-gray-900 dark:text-white">Available Groups</h2>
        <div class="grid gap-4 sm:grid-cols-2">
          <div 
            v-for="group in filteredGroups" 
            :key="group.id"
            class="bg-white dark:bg-gray-800 rounded-xl shadow-md hover:shadow-lg transition-all cursor-pointer p-4"
            @click="promptPassword(group)"
          >
            <div class="flex items-center space-x-3">
              <span class="text-2xl">🎯</span>
              <span class="font-semibold text-gray-900 dark:text-white">{{ group.name }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Create Group Form -->
      <div class="bg-gradient-to-r from-blue-500 to-indigo-500 rounded-xl shadow-lg p-6">
        <h2 class="text-2xl font-bold text-white mb-6">Create New Group</h2>
        <form @submit.prevent="createGroup" class="space-y-4">
          <input
            v-model="form.name"
            type="text"
            placeholder="Group Name"
            class="modern-input w-full bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
            required
          />
          <input
            v-model="form.password"
            type="password"
            placeholder="Password"
            class="modern-input w-full bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm"
            required
          />
          <button type="submit" class="modern-button w-full bg-gradient-to-r from-green-500 to-emerald-500 text-white font-bold">
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
    console.log('Raw API response:', response);
    
    if (Array.isArray(response.data)) {
      groups.value = response.data.map(group => ({
        id: group.id,
        name: group.name,
        created_at: group.created_at
      }));
      console.log('Processed groups:', groups.value);
    } else {
      console.error('Unexpected API response format:', response.data);
      error.value = 'Invalid data format received from server';
    }
  } catch (err) {
    console.error('Error fetching groups:', err);
    error.value = typeof err === 'string' ? err : 'Failed to fetch groups';
  } finally {
    loading.value = false;
  }
});

const promptPassword = async (group: Group) => {
  if (!group || !group.name) {
    console.error('Invalid group data in promptPassword:', group);
    return;
  }

  const password = prompt(`Enter password for group "${group.name}":`);
  if (!password) return;

  try {
    const response = await groupApi.getByName(group.name, password);
    const joinedGroup = response.data;
    groupStore.setGroup(joinedGroup, password);
    savedGroupsStore.addGroup(joinedGroup.name, joinedGroup.name, password);
    router.push(`/group/${joinedGroup.name}`);
  } catch (err) {
    console.error('Error joining group:', err);
    alert('Wrong password or group not found');
  }
};

const createGroup = async () => {
  try {
    const response = await groupApi.create(form.value.name, form.value.password);
    const group = response.data;
    groupStore.setGroup(group, form.value.password);
    savedGroupsStore.addGroup(group.name, group.name, form.value.password);
    router.push(`/group/${group.name}`);
    console.log('✅ Group created:', group);
  } catch (err) {
    console.error('Error creating group:', err);
    alert('Failed to create group. Please try again.');
  }
};

const joinSavedGroup = async (group: { name: string, password: string }) => {
  try {
    const response = await groupApi.getByName(group.name, group.password);
    const joinedGroup = response.data;
    groupStore.setGroup(joinedGroup, group.password);
    router.push(`/group/${group.name}`);
  } catch (err) {
    console.error('Error joining saved group:', err);
    alert('Failed to join group. The group may no longer exist.');
    savedGroupsStore.removeGroup(group.name);
  }
};

const removeGroup = (name: string) => {
  if (confirm('Are you sure you want to remove this group from your recent list?')) {
    savedGroupsStore.removeGroup(name);
  }
};
</script>
