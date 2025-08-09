import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Ref } from 'vue'

export const useTmaStore = defineStore('tma', () => {
  const token: Ref<unknown> = ref()

  function init() {
    token.value = window.Telegram?.WebApp?.initData || ''
    console.log('initData', token.value)
  }

  return { token, init }
})
