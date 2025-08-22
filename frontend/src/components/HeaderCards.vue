<script setup lang="ts">
import { ref } from 'vue'
import { useExpensesStore } from '@/stories/expensesStore'

import ExpenseStatistic from './ExpenseStatistic.vue';
import TotalAmountSmall from './TotalAmountSmall.vue';
import BudgetAmountSmall from './BudgetAmountSmall.vue';
import FilterButton from './FilterButton.vue';

const expensesStore = useExpensesStore()

const isCart = ref(false)
const isOpen = ref(true)

function handleClose() {
  isOpen.value = false
}

function handleOpen() {
  isOpen.value = true
}

function handleFilterChange() {
  if (isCart.value) {
    expensesStore.setFilter('paymentType', 'карта')
  } else {
    expensesStore.removeFilter('paymentType')
  }
}
</script>
<template>
  <header className="sticky top-0 py-8 px-4 bg-white-to-transparent">
    <!-- Filters -->
    <div class="flex flex-wrap gap-2 mb-6">
      <FilterButton label="Картой" v-model="isCart" @change="handleFilterChange" />
    </div>

    <ExpenseStatistic v-if="isOpen" :onClose="handleClose" />
    <div className="grid grid-cols-2 gap-6" v-else>
      <TotalAmountSmall @click="handleOpen" />
      <BudgetAmountSmall />
    </div>
  </header>
</template>

<style scoped>
.shadow-item {
  box-shadow: rgba(0, 0, 0, 0.12) 0 6px 34px 0;
}

.bg-white-to-transparent {
  background: linear-gradient(to bottom, white 80%, transparent);
}
</style>
