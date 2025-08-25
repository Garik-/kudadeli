<script setup lang="ts">
import {
  XMarkIcon
} from '@heroicons/vue/16/solid'

import { useExpensesStore } from '@/stories/expensesStore'
import { useFiltersStore } from '@/stories/filtersStore'

const expensesStore = useExpensesStore()
const filtersStore = useFiltersStore()

function setCategoryFilter(name: string) {
  filtersStore.setFilter('category', name)
}

const props = defineProps({
  onClose: { type: Function, required: true }
})

function handleClose() {
  props.onClose()
}

function getLabelColor(color: string, n: number) {
  const parts = color.split('-');
  parts[parts.length - 1] = n.toString();
  return parts.join('-');
}

</script>
<template>
  <div class="flex justify-between mb-6">
    <div>
      <div className="text-3xl font-bold">{{ expensesStore.totalAmount }}</div>
      <div className="text-sm">Траты</div>
    </div>
    <div className="w-8 h-8 bg-gray-50 rounded-full p-1 cursor-pointer" @click="handleClose">
      <XMarkIcon class="w-full h-full text-gray-500" />
    </div>
  </div>
  <div class="flex flex-wrap gap-2 mb-6">
    <div v-for="category in expensesStore.groupedByCategory" :key="category.amount"
      :class="['flex items-center gap-1 rounded-full p-1 cursor-pointer', getLabelColor(category.color, 100)]"
      @click="setCategoryFilter(category.name)">
      <div class="w-8 h-8 rounded-full p-2" :class="category.color">
        <component :is="category.icon" class="w-full h-full text-white" />
      </div>
      <div class="text-sm font-semibold text-gray-700 pr-2">{{ category.title }} {{ category.amountFormatted }}</div>
    </div>
  </div>
</template>
