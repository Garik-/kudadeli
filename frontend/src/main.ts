import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createMemoryHistory, createRouter } from 'vue-router'

import ExpenseList from '@/components/ExpenseList.vue'
import CategorySelector from './components/CategorySelector.vue'

import App from './App.vue'

const routes = [
  { path: '/', component: ExpenseList },
  { path: '/expenses/:id/category', name: 'edit-category', component: CategorySelector },
]

const router = createRouter({
  history: createMemoryHistory(),
  routes,
})

const pinia = createPinia()
const app = createApp(App)
app.use(pinia)
app.use(router)
app.mount('#app')
