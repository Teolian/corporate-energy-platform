import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/HomeView.vue')
  },
  {
    path: '/companies',
    name: 'companies',
    component: () => import('@/views/CompanyListView.vue')
  },
  {
    path: '/companies/:id',
    name: 'company-detail',
    component: () => import('@/views/CompanyDetailView.vue'),
    props: true
  },
  {
    path: '/compare',
    name: 'compare',
    component: () => import('@/views/CompareView.vue')
  }
]

const router = createRouter({
  history: createWebHistory('/'),
  routes
})

export default router
