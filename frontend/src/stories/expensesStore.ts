import { ref, computed } from 'vue'
import type { Expense } from '@/models/expense'
import { getTotalAmount } from '@/models/expense'

import { isToday, isYesterday } from 'date-fns'
import { fetchExpenses } from '@/services/api'
import { formatPrice, capitalizeFirstLetter } from '@/utils/formatter'
import { defineStore } from 'pinia'
import { useCategoriesStore } from './categoriesStore'
import { formatPercent } from '@/utils/formatter'

import type { FunctionalComponent } from 'vue'

export interface Item {
  id: string
  title: string
  category: string
  amount: string
  paymentType: string
}

export interface GroupedItems {
  date: string
  items: Item[]
}

export interface ExpenseByCategory {
  amount: number
  amountFormatted: string
  color: string
  percent: string
  name: string
  icon?: FunctionalComponent
}

export type GroupedAmount = Record<string, string>

export const useExpensesStore = defineStore('expenses', () => {
  const update = ref(true)

  const expenses = ref<Expense[]>([])
  const filter = ref<Partial<Expense>>({})

  const categoriesStore = useCategoriesStore()

  const filteredExpenses = computed(() => {
    return expenses.value.filter((expense) => {
      return Object.entries(filter.value).every(([key, value]) => {
        if (!value) return true // если фильтр по этому полю не задан
        return String(expense[key as keyof Expense]) === String(value)
      })
    })
  })

  const groupedTransactions = computed(() => transformExpenses(filteredExpenses.value))
  const groupedAmount = computed(() => transformExpensesAmount(filteredExpenses.value))

  const groupedByCategory = computed(() => transformExpensesByCategory(filteredExpenses.value))
  const totalAmount = computed(() => {
    const amount = getTotalAmount(filteredExpenses.value)
    return formatPrice(amount)
  })

  const loading = ref(false)
  const error = ref<string | null>(null)

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

  async function loadExpenses() {
    if (!update.value) {
      console.log('skip loadExpenses')
      return
    }

    loading.value = true
    error.value = null
    try {
      expenses.value = await fetchExpenses()
    } catch (e: unknown) {
      if (e instanceof Error) {
        error.value = e.message || 'Ошибка загрузки'
        console.error(error.value)
      } else {
        console.error('Unknown error', e)
      }
    } finally {
      loading.value = false
      update.value = false
    }
  }

  function needUpdate() {
    update.value = true
  }

  function transformExpensesByCategory(data: Expense[]): ExpenseByCategory[] {
    const c: Record<string, number> = {}
    let total = 0

    data.forEach((expense) => {
      const amount = parseFloat(expense.amount)
      total += amount

      if (!(expense.category in c)) {
        c[expense.category] = amount
      } else {
        c[expense.category] += amount
      }
    })

    const result = Object.entries(c).map(([name, amount]) => ({
      amount,
      amountFormatted: formatPrice(amount),
      color: categoriesStore.getColorByName(name),
      name: capitalizeFirstLetter(name),
      icon: categoriesStore.getIconComponentByName(name),
      percent: formatPercent((amount / total) * 100),
    }))

    result.sort((a, b) => b.amount - a.amount)

    return result
  }

  function formatDateToGroupLabel(dateString: string) {
    const date = new Date(dateString)

    if (isToday(date)) return 'Сегодня'
    if (isYesterday(date)) return 'Вчера'

    return date.toLocaleDateString('ru-RU', {
      day: 'numeric',
      month: 'long',
    }) // например: 16 июля
  }

  function transformExpensesAmount(data: Expense[]): Record<string, string> {
    const amounts: Record<string, number> = {}

    data.forEach((expense) => {
      const groupKey = formatDateToGroupLabel(expense.createdAt)
      const amount = parseFloat(expense.amount)

      if (!amounts[groupKey]) {
        amounts[groupKey] = 0
      }

      amounts[groupKey] += amount
    })

    // Преобразуем объект amounts в объект с отформатированными значениями
    const formattedAmounts: Record<string, string> = {}
    for (const key in amounts) {
      formattedAmounts[key] = formatPrice(amounts[key]) // форматируем каждое значение
    }
    return formattedAmounts
  }

  function transformExpenses(data: Expense[]): GroupedItems[] {
    const map: Map<string, Item[]> = new Map()

    data.forEach((expense) => {
      const groupKey = formatDateToGroupLabel(expense.createdAt)

      if (!map.has(groupKey)) {
        map.set(groupKey, [])
      }

      map.get(groupKey)!.push({
        id: expense.id,
        title: capitalizeFirstLetter(expense.description),
        category: expense.category,
        amount: formatPrice(parseFloat(expense.amount)),
        paymentType: expense.paymentType,
        // logo: getLogo(expense), // логика для логотипа
      })
    })

    // Преобразуем Map в массив
    return Array.from(map.entries()).map(([date, items]) => ({
      date,
      items,
    }))
  }

  return {
    groupedTransactions,
    groupedAmount,
    expenses,
    loading,
    error,
    loadExpenses,
    needUpdate,
    filteredExpenses,
    setFilter,
    removeFilter,
    clearFilter,
    groupedByCategory,
    totalAmount,
  }
})
