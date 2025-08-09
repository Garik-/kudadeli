import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Ref } from 'vue'
import { retrieveLaunchParams } from '@telegram-apps/sdk'

export const useTmaStore = defineStore('tma', () => {
  const token: Ref<unknown> = ref()

  function init() {
    try {
      const { initDataRaw, initData } = retrieveLaunchParams()
      token.value = initData

      console.log(initData, initDataRaw)
    } catch (e: unknown) {
      if (e instanceof Error) {
        console.log('tmaStore init', e.message)
      }
    }
  }

  return { token, init }
})
