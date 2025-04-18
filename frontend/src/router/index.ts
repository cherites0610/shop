import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/view/HomeView.vue'
import CommodityView from '@/view/CommodityView.vue'
import LineLogin from '@/view/LineLogin.vue'
import SpecificationView from '@/view/SpecificationView.vue'

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
    {
      path:'/LineLogin',
      name:"LineLogin",
      component: LineLogin
    },
    {
      path:"/specification/:id",
      name:"specification",
      component: SpecificationView
    }
  ],
})

export default router
