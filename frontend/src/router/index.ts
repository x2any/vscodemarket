import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'home', component: () => import('../pages/HomePage.vue') },
  { path: '/search', name: 'search', component: () => import('../pages/ExtensionSearch.vue') },
  { path: '/extension/:pub/:name', name: 'extension-detail', component: () => import('../pages/ExtensionDetail.vue') },
  { path: '/extension/:pub/:name/v/:ver', name: 'extension-version', component: () => import('../pages/ExtensionDetail.vue') },
  { path: '/releases', name: 'releases', component: () => import('../pages/VersionList.vue') },
  { path: '/trending', name: 'trending', component: () => import('../pages/Trending.vue') },
  { path: '/docs', name: 'docs', component: () => import('../pages/Docs.vue') }
]

export const router = createRouter({
  history: createWebHistory(),
  routes
})
