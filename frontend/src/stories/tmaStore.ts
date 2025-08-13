import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Ref } from 'vue'

import type { Telegram } from 'telegram-web-app'

declare global {
  interface Window {
    Telegram: Telegram
  }
}

export const useTmaStore = defineStore('tma', () => {
  const token: Ref<unknown> = ref()

  function init() {
    token.value = window.Telegram?.WebApp?.initData || ''
    console.log('initData', token.value)
  }

  return { token, init }
})
