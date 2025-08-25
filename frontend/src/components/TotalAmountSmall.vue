<script setup lang="ts">
import { useExpensesStore } from '@/stories/expensesStore'
const expenseStore = useExpensesStore()
</script>
<template>
  <div className="flex flex-col bg-white p-6 rounded-2xl shadow-item">
    <div className="font-bold text-lg">{{ expenseStore.totalAmount }}</div>
    <div className="text-gray-500 text-sm mb-4">Траты</div>

    <div className="flex h-3 w-full rounded-full overflow-hidden">
      <div
        v-for="group in expenseStore.groupedByCategory"
        :key="group.amount"
        :class="group.color"
        :style="{ width: group.percent }"
      ></div>
    </div>
  </div>
</template>
