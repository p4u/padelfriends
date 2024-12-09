<template>
  <div class="space-y-4">
    <h2 class="text-xl font-semibold text-gray-900 dark:text-white">🎾 Create New Group</h2>
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <input
        v-model="form.name"
        type="text"
        placeholder="Group Name"
        class="modern-input w-full text-gray-900 dark:text-white"
        required
      />
      <input
        v-model="form.password"
        type="password"
        placeholder="Password"
        class="modern-input w-full text-gray-900 dark:text-white"
        required
      />
      <button type="submit" class="modern-button w-full">
        Create Group
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const form = ref({ name: '', password: '' });

const emit = defineEmits<{
  (e: 'create', name: string, password: string): void;
}>();

const handleSubmit = () => {
  emit('create', form.value.name, form.value.password);
  form.value = { name: '', password: '' };
};
</script>
