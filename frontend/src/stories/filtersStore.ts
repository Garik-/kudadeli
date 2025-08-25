import type { Expense } from '@/models/expense'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export type Filter = Partial<Expense>

export const useFiltersStore = defineStore('filters', () => {
  const filter = ref<Filter>({})

  function setFilter<K extends keyof Expense>(key: K, value: Expense[K] | null) {
    if (value == null || value === '') {
      delete filter.value[key]
    } else {
      filter.value[key] = value
    }
  }

  function removeFilter<K extends keyof Expense>(key: K) {
    delete filter.value[key]
  }

  function clearFilter() {
    filter.value = {}
  }

  return {
    filter,
    setFilter,
    removeFilter,
    clearFilter,
  }
})
