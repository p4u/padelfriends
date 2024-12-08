import { createRouter, createWebHistory } from 'vue-router';
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
      path: '/group/:id',
      name: 'group',
      component: Group,
      props: true,
    },
  ],
});