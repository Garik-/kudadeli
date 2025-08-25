<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

import { useExpensesStore } from '@/stories/expensesStore'
import { useCategoriesStore } from '@/stories/categoriesStore'

import FullScreenLoader from '@/components/FullScreenLoader.vue'
import HeaderCards from './HeaderCards.vue'

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

<template>
  <FullScreenLoader v-if="expensesStore.loading"></FullScreenLoader>

  <template v-else>
    <div class="max-w-3xl mx-auto">
      <HeaderCards />

      <main className="space-y-8 pb-4">
        <!-- Transactions by Date -->
        <div v-for="group in expensesStore.groupedTransactions" :key="group.date" class="space-y-2">
          <div class="flex items-center justify-between px-4">
            <div class="text-lg font-bold">{{ group.date }}</div>
            <div class="text-right text-gray-400">
              -{{ expensesStore.groupedAmount[group.date] }}
            </div>
          </div>

          <div class="space-y-2">
            <div v-for="tx in group.items" :key="tx.id">
              <div
                @click="selectItem(tx.id, tx.category)"
                class="flex justify-between items-center active:bg-gray-50 px-4 py-2 cursor-pointer"
              >
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
