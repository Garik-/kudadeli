<script setup lang="ts">
import { onMounted } from 'vue'
import FullScreenLoader from '@/components/FullScreenLoader.vue'
import { useExpensesStore } from '@/stories/expensesStore'

import { useCategoriesStore } from '@/stories/categoriesStore';
import { useRouter } from 'vue-router';
import TotalAmountSmall from './TotalAmountSmall.vue';
import BudgetAmountSmall from './BudgetAmountSmall.vue';

import { XMarkIcon } from '@heroicons/vue/16/solid'

const store = useCategoriesStore()
const router = useRouter()
const expensesStore = useExpensesStore()

onMounted(() => {
  expensesStore.loadExpenses()
})

function selectItem(ID: string, category: string) {
  router.push({ name: 'edit-category', params: { id: ID + ':' + store.getIdByName(category) } })
}
</script>


<style scoped>
.shadow-item {
  box-shadow: rgba(0, 0, 0, 0.12) 0 6px 34px 0;
}

.filter-active {
  padding-left: 12px;
  padding-right: 6px;
}

.filter {
  padding-left: 12px;
  padding-right: 12px;
}
</style>

<template>

  <FullScreenLoader v-if="expensesStore.loading"></FullScreenLoader>

  <template v-else>
    <div class="max-w-3xl mx-auto">

      <header className="sticky top-0 py-8 px-4">
        <!-- Filters -->
        <div class="flex flex-wrap gap-2 mb-6">
          <button
            class="filter-active inline-flex items-center justify-between gap-1 py-2 rounded-2xl bg-blue-500 text-white text-sm font-semibold cursor-pointer"><span>Black</span>
            <XMarkIcon class="w-[1.5em] h-[1.5em]" />

          </button>
          <button class="filter py-2 rounded-2xl bg-gray-100 text-gray-700 text-sm font-bold cursor-pointer">Без
            переводов</button>
        </div>
        <div className="grid  grid-cols-2 gap-6 rounded-2xl">

          <TotalAmountSmall />
          <BudgetAmountSmall />
        </div>
      </header>

      <main className="space-y-8 pb-4">


        <!-- Transactions by Date -->
        <div v-for="group in expensesStore.groupedTransactions" :key="group.date" class="space-y-2">
          <div class="flex items-center justify-between px-4">
            <div class="text-lg font-bold">{{ group.date }}</div>
            <div class="text-right text-gray-400">-{{ expensesStore.groupedAmount[group.date] }}</div>
          </div>

          <div class="space-y-2">
            <div v-for="tx in group.items" :key="tx.id">
              <div @click="selectItem(tx.id, tx.category)"
                class="flex justify-between items-center active:bg-gray-50 px-4 py-2 cursor-pointer">
                <div class="flex items-center gap-4">
                  <div>
                    <div class="font-medium">{{ tx.title }}</div>
                    <div class="text-sm text-gray-500">{{ tx.category }}</div>
                  </div>
                </div>
                <div class="text-right">
                  <div class="text-red-600 font-medium">-{{ tx.amount }}</div>
                  <div class="text-sm text-gray-500">{{ tx.paymentType }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </template>
</template>
