<template>
  <div class="max-w-md mx-auto space-y-8">
    <div class="modern-container bg-white dark:bg-gray-800">
      <h1 class="text-2xl font-bold text-center text-gray-900 dark:text-white mb-8">
        🎾 Padel Friends
      </h1>
      
      <!-- Recent Groups -->
      <div v-if="savedGroups.length > 0" class="mb-8">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">🏆 Recent Groups</h2>
        <SavedGroupList
          :groups="savedGroups"
          @select="joinSavedGroup"
          @remove="removeGroup"
        />
      </div>

      <!-- Available Groups -->
      <div v-if="filteredGroups.length > 0" class="mb-8">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">🌟 Available Groups</h2>
        <GroupList
          :groups="filteredGroups"
          :loading="loading"
          :error="error"
          @select="promptPassword"
        />
      </div>

      <!-- Create Group Form -->
      <CreateGroupForm @create="createGroup" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import { groupApi } from '../api';
import { useGroupStore } from '../stores/group';
import { useSavedGroupsStore } from '../stores/savedGroups';
import CreateGroupForm from '../components/CreateGroupForm.vue';
import GroupList from '../components/GroupList.vue';
import SavedGroupList from '../components/SavedGroupList.vue';
import type { Group } from '../types';

const router = useRouter();
const groupStore = useGroupStore();
const savedGroupsStore = useSavedGroupsStore();

const groups = ref<Group[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);

const { savedGroups } = storeToRefs(savedGroupsStore);

// First ensure groups data is valid
const processedGroups = computed(() => {
  return groups.value.filter(group => {
    const isValid = group && typeof group === 'object' && 'id' in group && 'name' in group;
    if (!isValid) {
      console.error('Invalid group data found:', group);
    }
    return isValid;
  });
});

// Then filter out groups that are already in savedGroups
const filteredGroups = computed(() => {
  const savedGroupNames = new Set(savedGroups.value.map(group => group.name));
  return processedGroups.value.filter(group => !savedGroupNames.has(group.name));
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

  const password = prompt(`🔒 Enter password for group "${group.name}":`);
  if (!password) return;

  try {
    const response = await groupApi.getByName(group.name, password);
    const joinedGroup = response.data;
    groupStore.setGroup(joinedGroup, password);
    savedGroupsStore.addGroup(joinedGroup.name, joinedGroup.name, password);
    router.push(`/group/${joinedGroup.name}`);
  } catch (err) {
    console.error('Error joining group:', err);
    alert('❌ Wrong password or group not found');
  }
};

const createGroup = async (name: string, password: string) => {
  try {
    const response = await groupApi.create(name, password);
    const group = response.data;
    groupStore.setGroup(group, password);
    savedGroupsStore.addGroup(group.name, group.name, password);
    router.push(`/group/${group.name}`);
    console.log('✅ Group created:', group);
  } catch (err) {
    console.error('Error creating group:', err);
    alert('❌ Failed to create group. Please try again.');
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
    alert('❌ Failed to join group. The group may no longer exist.');
    savedGroupsStore.removeGroup(group.name);
  }
};

const removeGroup = (name: string) => {
  if (confirm('Are you sure you want to remove this group from your recent list?')) {
    savedGroupsStore.removeGroup(name);
  }
};
</script>
