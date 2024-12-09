import { createRouter, createWebHistory } from 'vue-router';
import { useGroupStore } from '../stores/group';
import Home from '../views/Home.vue';
import Group from '../views/Group.vue';

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/group/:name',
      name: 'group',
      component: Group,
      props: true,
      beforeEnter: async (to, from, next) => {
        const groupStore = useGroupStore();
        const groupName = to.params.name as string;

        // If we're already in the correct group, proceed
        if (groupStore.currentGroup?.name === groupName) {
          next();
          return;
        }

        // Try to restore group from storage
        const restored = await groupStore.restoreGroupFromStorage();
        if (restored && groupStore.currentGroup?.name === groupName) {
          next();
        } else {
          // If restoration fails or it's a different group, redirect to home
          next('/');
        }
      }
    },
  ],
});

// Global navigation guard to handle initial page load
router.beforeEach(async (to, from, next) => {
  // Skip for home route
  if (to.name === 'home') {
    next();
    return;
  }

  const groupStore = useGroupStore();
  
  // If we already have a group loaded, continue
  if (groupStore.currentGroup) {
    next();
    return;
  }

  // Try to restore group from storage
  const restored = await groupStore.restoreGroupFromStorage();
  if (restored) {
    next();
  } else if (to.name !== 'home') {
    // If restoration fails and we're not going to home, redirect to home
    next('/');
  } else {
    next();
  }
});

export default router;
