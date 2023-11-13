import Vue from 'vue'
import VueRouter from 'vue-router'
import Login from '../components/Login.vue'
import Home from '../components/Home.vue'
import Welcome from '../components/Welcome.vue'
import Adminer from '../components/platform/Adminer.vue'
import AdminerGroup from '../components/platform/AdminerGroup.vue'
import AdminerAuth from '../components/platform/AdminerAuth.vue'

Vue.use(VueRouter)

const routes = [
  { path: '/', redirect: '/home' },
  { path: '/login', component: Login },
  {
    path: '/home',
    component: Home,
    redirect: '/welcome',
    children: [
      { path: '/welcome', component: Welcome },
      { path: '/adminer/list', component: Adminer },
      { path: '/adminer/group/list', component: AdminerGroup },
      { path: '/adminer/auth/list', component: AdminerAuth }
    ]
  }
]

const router = new VueRouter({
  routes
})

// 挂载路由守卫导航
// eslint-disable-next-line no-use-before-define
router.beforeEach((to, from, next) => {
  // 如果直接访问登录页面直接放行
  if (to.path === '/login') {
    next()
  }

  // 如果访问其他页面先获取token,获取不到直接跳转登录页面
  const token = localStorage.getItem('token')
  if (!token) {
    next('/login')
  }

  next()
})

export default router
