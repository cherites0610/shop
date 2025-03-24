import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/view/HomeView.vue'
import App from '@/App.vue'
import CommodityView from '@/view/CommodityView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/commodity/:id',
      name: 'commodity',
      component: CommodityView,
    },
  ],
})

export default router
