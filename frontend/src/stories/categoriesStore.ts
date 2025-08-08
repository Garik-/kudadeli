import { ref } from 'vue'
import type { Ref } from 'vue'
import type { Category } from '@/models/category'
import { fetchCategories } from '@/services/api'
import { capitalizeFirstLetter } from '@/utils/formatter'

import { defineStore } from 'pinia'

export const useCategoriesStore = defineStore('categories', () => {
  const categories: Ref<Category[]> = ref([])
  const categoriesMap: Ref<Record<string, number>> = ref({})
  const loading = ref(false)
  const error = ref<string | null>(null)

  function getIdByName(name: string): number {
    if (!(name in categoriesMap.value)) {
      return 0
    }

    return categoriesMap.value[name]
  }

  async function loadCategories() {
    loading.value = true
    error.value = null

    try {
      const data = await fetchCategories()
      const c: Category[] = []
      data.forEach(({ id, name }) => {
        c.push({ id, name: capitalizeFirstLetter(name) })
        categoriesMap.value[name] = id
      })
      categories.value = c
    } catch (e: unknown) {
      if (e instanceof Error) {
        error.value = e.message || 'Ошибка загрузки'
        console.error(error.value)
      } else {
        console.error('Unknown error', e)
      }
    } finally {
      loading.value = false
    }
  }

  return {
    categories,
    loading,
    error,
    loadCategories,
    getIdByName,
  }
})
