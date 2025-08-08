<script setup lang="ts">
import { onMounted } from 'vue'
import FullScreenLoader from '@/components/FullScreenLoader.vue'
import { useExpenses } from '@/composables/useExpenses'

import { useCategoriesStore } from '@/stories/categoriesStore';

const store = useCategoriesStore()


const { groupedTransactions,
  groupedAmount,
  totalAmount,
  budgetAmount,
  loading, loadExpenses } = useExpenses()

onMounted(() => {
  loadExpenses()
})

function selectItem(ID: string, category: string) {
  window.location.hash = encodeURIComponent(ID + ':' + store.getIdByName(category))
}
</script>


<style scoped>
.shadow-item {
  box-shadow: rgba(0, 0, 0, 0.12) 0 6px 34px 0;
}
</style>

<template>

  <FullScreenLoader v-if="loading"></FullScreenLoader>

  <template v-else>
    <div class="max-w-3xl mx-auto">

      <div className="sticky top-0 py-8 px-4">
        <div className="grid  grid-cols-2 gap-4 rounded-2xl">

          <div className="flex flex-col bg-white p-6 rounded-2xl shadow-item">
            <div className="font-bold text-lg">{{ totalAmount }}</div>
            <div className="text-gray-500 text-sm">Траты</div>
            <!--<div className="flex h-4 w-full rounded-full overflow-hidden">

          <div className="bg-blue-400 flex-grow"></div>

          <div className="bg-indigo-300 w-6"></div>
          <div className="bg-rose-500 w-6"></div>
          <div className="bg-pink-500 w-6"></div>
          <div className="bg-amber-400 w-6"></div>
          <div className="bg-slate-400 w-6"></div>
        </div> -->
          </div>

          <div className="flex flex-col bg-white p-6 rounded-2xl shadow-item">
            <div className="font-bold text-lg">{{ budgetAmount }}</div>
            <div className="text-gray-500 text-sm">Бюджет</div>
            <!--<div className="flex h-4 w-full rounded-full overflow-hidden">
          <div className="bg-blue-400 flex-grow"></div>
          <div className="bg-sky-600 w-6"></div>
          <div className="bg-teal-300 w-6"></div>
        </div>-->
          </div>
        </div>
      </div>

      <div className="space-y-8 relative">

        <!-- Transactions by Date -->
        <div v-for="group in groupedTransactions" :key="group.date" class="space-y-2">
          <div class="flex items-center justify-between px-4">
            <div class="text-lg font-bold">{{ group.date }}</div>
            <div class="text-right text-gray-400">-{{ groupedAmount[group.date] }}</div>
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

      </div>
    </div>
  </template>
</template>
