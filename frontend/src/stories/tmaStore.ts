import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Ref } from 'vue'

export const useTmaStore = defineStore('tma', () => {
  const token: Ref<unknown> = ref()

  function init() {
    try {
      const { initData } = window.Telegram?.WebApp?.initData || undefined
      token.value = initData
      console.log(initData)
    } catch (e: unknown) {
      if (e instanceof Error) {
        console.log('tmaStore init', e.message)
      }
    }
  }

  return { token, init }
})
