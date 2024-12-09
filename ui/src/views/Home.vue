<template>
  <div class="max-w-4xl mx-auto space-y-8">
    <!-- Create Group Form -->
    <CreateGroupForm @create-group="handleCreateGroup" />

    <!-- Available Groups -->
    <div class="modern-container bg-white dark:bg-gray-800">
      <h2 class="text-xl font-bold mb-4 text-gray-900 dark:text-white">
        {{ t('groups.availableGroups') }}
      </h2>
      <div v-if="loading" class="text-center text-gray-500 dark:text-gray-400">
        {{ t('common.loading') }}
      </div>
      <div v-else-if="error" class="text-center text-red-500">
        {{ error }}
      </div>
      <div v-else>
        <GroupList 
          :groups="groups"
          @join-group="handleJoinGroup"
        />
      </div>
    </div>

    <!-- Recent Groups -->
    <div class="modern-container bg-white dark:bg-gray-800">
      <h2 class="text-xl font-bold mb-4 text-gray-900 dark:text-white">
        {{ t('groups.recentGroups') }}
      </h2>
      <SavedGroupList @join-group="handleJoinGroup" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from '../i18n';
import { groupApi } from '../api';
import { useGroupStore } from '../stores/group';
import { CreateGroupForm, GroupList, SavedGroupList } from '../components';
import type { Group } from '../types';

const router = useRouter();
const groupStore = useGroupStore();
const { t } = useI18n();

const groups = ref<Group[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);

const loadGroups = async () => {
  loading.value = true;
  error.value = null;

  try {
    const response = await groupApi.getGroups();
    groups.value = response.data;
  } catch (err) {
    console.error('Error fetching groups:', err);
    error.value = t('errors.failedToLoad');
  } finally {
    loading.value = false;
  }
};

const handleCreateGroup = async (name: string, password: string) => {
  try {
    await groupApi.createGroup(name, password);
    await groupStore.setGroup(name, password);
    router.push(`/group/${name}`);
  } catch (error) {
    console.error('Error creating group:', error);
    alert(t('errors.failedToCreate'));
  }
};

const handleJoinGroup = async (name: string, password: string) => {
  try {
    await groupStore.setGroup(name, password);
    router.push(`/group/${name}`);
  } catch (error) {
    console.error('Error joining group:', error);
    alert(t('errors.invalidPassword'));
  }
};

onMounted(() => {
  loadGroups();
});
</script>
