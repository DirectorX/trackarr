import Vue from 'vue'
import VueRouter from 'vue-router'
const Home  = () => import('../views/Home.vue')
const Logs = () => import('../views/Logs.vue')
const Status = () => import('../views/Status.vue')

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home
  },
  {
    path: '/logs',
    name: 'logs',
    component: Logs
  },
  {
    path: '/status',
    name: 'status',
    component: Status
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
