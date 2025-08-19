<script setup lang="ts">
import { computed } from 'vue'
import { getTotalAmount } from '@/models/expense'
import { useExpensesStore } from '@/stories/expensesStore'
import { formatPrice, formatPercent } from '@/utils/formatter'

const BUDGET = 3_000_000.0
const expenseStore = useExpensesStore()
const budgetAmount = computed(() => {
  const amount = getTotalAmount(expenseStore.expenses)
  return formatPrice(BUDGET - amount)
})
const budgetPercent = computed(() => {
  const amount = getTotalAmount(expenseStore.expenses)
  return formatPercent(100 - (amount / BUDGET) * 100)
})
</script>
<template>
  <div className="flex flex-col bg-white p-6 rounded-2xl shadow-item">
    <div className="font-bold text-lg">{{ budgetAmount }}</div>
    <div className="text-gray-500 text-sm mb-4">Бюджет</div>
    <div className="flex h-3 w-full bg-gray-200 rounded-full overflow-hidden">
      <div className="bg-blue-400 rounded-full" :style="{ width: budgetPercent }"></div>
    </div>
  </div>
</template>
