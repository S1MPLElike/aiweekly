import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import Login from '../views/Login.vue'
import Dashboard from '../views/Dashboard.vue'
import { getStoredUser } from '../services/auth'
import { initAuth } from '../main'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

let authInitialized = false

router.beforeEach(async (to, from) => {
  if (!authInitialized) {
    await initAuth()
    authInitialized = true
  }
  
  const user = getStoredUser()
  
  if (to.path === '/login' && user) {
    return '/dashboard'
  }
  
  if (to.path === '/dashboard' && !user) {
    return '/login'
  }
})

export default router
