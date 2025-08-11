import { ref } from 'vue'
import type { Expense } from '@/models/expense'
import { isToday, isYesterday } from 'date-fns'
import { fetchExpenses } from '@/services/api'
import { formatPrice, capitalizeFirstLetter, formatPercent } from '@/utils/formatter'
import { defineStore } from 'pinia'

const BUDGET = 3_000_000.0

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
}

export type GroupedAmount = Record<string, string>

export const useExpensesStore = defineStore('expenses', () => {
  const update = ref(true)
  const groupedTransactions = ref<GroupedItems[]>([])
  const groupedAmount = ref<GroupedAmount>({})
  const groupedByCategory = ref<ExpenseByCategory[]>([])

  const totalAmount = ref('')
  const budgetAmount = ref('')
  const budgetPercent = ref('')

  const loading = ref(false)
  const error = ref<string | null>(null)

  async function loadExpenses() {
    if (!update.value) {
      console.log('skip loadExpenses')
      return
    }

    loading.value = true
    error.value = null
    try {
      const data = await fetchExpenses()
      groupedTransactions.value = transformExpenses(data)
      groupedAmount.value = transformExpensesAmount(data)
      groupedByCategory.value = transformExpensesByCategory(data)

      const amount = getTotalAmount(data)

      totalAmount.value = formatPrice(amount)
      budgetAmount.value = formatPrice(BUDGET - amount)
      budgetPercent.value = formatPercent(100 - (amount / BUDGET) * 100)
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

  return {
    groupedTransactions,
    groupedAmount,
    groupedByCategory,
    totalAmount,
    budgetAmount,
    budgetPercent,
    loading,
    error,
    loadExpenses,
    needUpdate,
  }
})

function getTotalAmount(data: Expense[]) {
  let amount = 0

  data.forEach((expense) => {
    amount += parseFloat(expense.amount)
  })

  return amount
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

  const colors = [
    'bg-indigo-300',
    'bg-rose-500',
    'bg-pink-500',
    'bg-amber-400',
    'bg-slate-400',
    'bg-blue-400',
  ] // TODO: кароче надо сделать все таки в сторе категорий структуру [id] = {name, color} - или пофиг или типа {id, name, color, icon}[]
  // и уже исходя из этого строить цвета и проценты - потому что без этого компонент круга не сделаешь

  const result = Object.values(c).map((amount) => ({
    amount,
    amountFormatted: formatPrice(amount),
    color: colors.pop() || 'bg-gray-400',
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
