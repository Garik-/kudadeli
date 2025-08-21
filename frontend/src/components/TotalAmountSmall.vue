<script setup lang="ts">
import type { Expense } from '@/models/expense'

import { computed } from 'vue'
import { useExpensesStore } from '@/stories/expensesStore'
import { useCategoriesStore } from '@/stories/categoriesStore'
import { formatPrice, formatPercent } from '@/utils/formatter'
import { getTotalAmount } from '@/models/expense'

interface ExpenseByCategory {
  amount: number
  amountFormatted: string
  color: string
  percent: string
}

const categoriesStore = useCategoriesStore()

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
