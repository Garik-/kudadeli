import { ref } from 'vue'
import type { FunctionalComponent, Ref } from 'vue'
import { fetchCategories } from '@/services/api'
import { capitalizeFirstLetter } from '@/utils/formatter'
import {
  ArchiveBoxIcon,
  CurrencyDollarIcon,
  WrenchScrewdriverIcon,
  HomeModernIcon,
  QuestionMarkCircleIcon,
} from '@heroicons/vue/24/solid'

import { defineStore } from 'pinia'

export type CategoryID = number
export interface CategoryItem {
  id: CategoryID
  name: string
  color: string
  icon?: FunctionalComponent
}

export const useCategoriesStore = defineStore('categories', () => {
  const categories: Ref<CategoryItem[]> = ref([])
  const categoriesMap: Ref<Record<string, CategoryID>> = ref({})
  const loading = ref(false)
  const error = ref<string | null>(null)

  const colors: Record<number, string> = {
    1: 'bg-indigo-300',
    2: 'bg-rose-500',
    3: 'bg-pink-500',
    4: 'bg-amber-400',
    5: 'bg-slate-400',
    6: 'bg-blue-400',
  }

  const colorsHex: Record<number, string> = {
    'bg-indigo-300': 'oklch(78.5% 0.115 274.713)',
    'bg-rose-500': 'oklch(64.5% 0.246 16.439)',
    'bg-pink-500': 'oklch(65.6% 0.241 354.308)',
    'bg-amber-400': 'oklch(82.8% 0.189 84.429)',
    'bg-slate-400': 'oklch(70.4% 0.04 256.788)',
    'bg-blue-400': 'oklch(70.7% 0.165 254.624)',
  }

  const icons: Record<number, FunctionalComponent> = {
    1: ArchiveBoxIcon,
    2: CurrencyDollarIcon,
    3: WrenchScrewdriverIcon,
    4: HomeModernIcon,
    5: QuestionMarkCircleIcon,
  }

  function getIconComponent(id: CategoryID) {
    if (id in icons) {
      return icons[id]
    }
    return undefined
  }

  function getColorById(id: CategoryID) {
    if (!(id in colors)) {
      return 'bg-gray-400'
    }

    return colors[id]
  }

  function getColorByName(name: string) {
    return getColorById(getIdByName(name))
  }

  function getIdByName(name: string): CategoryID {
    if (!(name in categoriesMap.value)) {
      return 0
    }

    return categoriesMap.value[name]
  }

  function getIconComponentByName(name: string) {
    return getIconComponent(getIdByName(name))
  }

  function getHexColor(name: string) {
    if (!(name in colorsHex)) {
      return 'oklch(70.7% 0.022 261.325)'
    }
    return colorsHex[name]
  }

  async function loadCategories() {
    loading.value = true
    error.value = null

    try {
      const data = await fetchCategories()
      const c: CategoryItem[] = []
      data.forEach(({ id, name }) => {
        c.push({
          id,
          name: capitalizeFirstLetter(name),
          color: getColorById(id),
          icon: getIconComponent(id),
        })
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
    getColorByName,
    getIconComponentByName,
    getHexColor,
  }
})
