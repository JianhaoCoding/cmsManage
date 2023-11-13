import Vue from 'vue'
import App from './App.vue'
import router from './router'
import './plugins/element.js'
import './assets/css/global.css'
import axios from 'axios'
axios.defaults.baseURL = 'http://127.0.0.1:8082/'
axios.interceptors.request.use(config => {
  const itemStr = localStorage.getItem('token')
  let token = ''
  if (itemStr) {
    const item = JSON.parse(itemStr)
    const now = new Date()
    // 检查是否过期
    if (now.getTime() > item.expiry) {
      // 如果过期，从 localStorage 中删除该项
      localStorage.removeItem('token')
    } else {
      token = item.value
    }
  }

  config.headers.Authorization = token
  return config
})
// axios.defaults.headers.post['Content-Type'] = 'application/json'
Vue.prototype.$http = axios

// 写入localStorage缓存
Vue.prototype.$setLocalStorage = function (key, value, ttl) {
  const now = new Date()
  const item = {
    value: value,
    expiry: now.getTime() + ttl
  }
  localStorage.setItem(key, JSON.stringify(item))
}

// 读取localStorage缓存
Vue.prototype.$getLocalStorage = function (key) {
  const itemStr = localStorage.getItem(key)

  // 如果不存在，则返回 null
  if (!itemStr) {
    return null
  }

  const item = JSON.parse(itemStr)
  const now = new Date()

  // 检查是否过期
  if (now.getTime() > item.expiry) {
    // 如果过期，从 localStorage 中删除该项
    localStorage.removeItem(key)
    return null
  }
  return item.value
}

Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
