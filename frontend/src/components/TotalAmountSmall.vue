<script setup lang="ts">
import type { Expense } from '@/models/expense'

import { computed } from 'vue'
import { useExpensesStore } from '@/stories/expensesStore'
import { formatPrice, formatPercent } from '@/utils/formatter'
import { getTotalAmount } from '@/models/expense'

interface ExpenseByCategory {
  amount: number
  amountFormatted: string
  color: string
  percent: string
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

const expenseStore = useExpensesStore()

// TODO: надо перенести в компонент он не всегда показывается
const groupedByCategory = computed(() => transformExpensesByCategory(expenseStore.filteredExpenses))

// TODO: надо перенести в компонент он не всегда показывается
const totalAmount = computed(() => {
  const amount = getTotalAmount(expenseStore.filteredExpenses)
  return formatPrice(amount)
})
</script>
<template>
  <div className="flex flex-col bg-white p-6 rounded-2xl shadow-item">
    <div className="font-bold text-lg">{{ totalAmount }}</div>
    <div className="text-gray-500 text-sm mb-4">Траты</div>

    <div className="flex h-3 w-full rounded-full overflow-hidden">
      <div v-for="group in groupedByCategory" :key="group.amount" :class="group.color"
        :style="{ width: group.percent }">
      </div>
    </div>
  </div>
</template>
