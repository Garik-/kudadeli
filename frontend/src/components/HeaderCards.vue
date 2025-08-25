<script setup lang="ts">
import { ref } from 'vue'

import ExpenseStatistic from './ExpenseStatistic.vue';
import TotalAmountSmall from './TotalAmountSmall.vue';
import BudgetAmountSmall from './BudgetAmountSmall.vue';
import FiltersList from './FiltersList.vue';

const isOpen = ref(false)

function handleClose() {
  isOpen.value = false
}

function handleOpen() {
  isOpen.value = true
}


</script>
<template>
  <header className="sticky top-0 py-8 px-4 bg-white-to-transparent">
    <!-- Filters -->
    <FiltersList />

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
