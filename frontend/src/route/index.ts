import LoginPage from '../components/Login/index.vue'
import { createRouter, createWebHistory } from 'vue-router'
const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: LoginPage },
]
const router = createRouter({
  history: createWebHistory(),
  routes,
})
