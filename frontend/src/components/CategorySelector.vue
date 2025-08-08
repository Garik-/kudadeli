<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RadioGroup, RadioGroupOption } from '@headlessui/vue'
import FullScreenLoader from '@/components/FullScreenLoader.vue'


import {
  CheckCircleIcon,
  ArrowLeftIcon,
} from '@heroicons/vue/24/solid'

import { useCategoriesStore } from '@/stories/categoriesStore'




const selected = ref(0)
const id = ref('')

const store = useCategoriesStore()


function goBack() {
  window.history.back()
}

function changeCategory() {
  console.log(selected.value)
}

onMounted(() => {
  store.loadCategories()

  const [itemID, itemCategoryID] = decodeURIComponent(window.location.hash.substring(1)).split(":");
  id.value = itemID
  selected.value = parseInt(itemCategoryID)
})

</script>

<template>
  <FullScreenLoader v-if="store.loading"></FullScreenLoader>
  <template v-else>
    <div class="flex flex-col min-h-screen bg-white text-black">
      <!-- Header -->
      <header class="sticky top-0 bg-white z-10 px-4 pt-4 pb-2">
        <button @click="goBack" class="text-blue-400">
          <ArrowLeftIcon class="w-6 h-6" />
        </button>
        <h1 class="text-2xl font-bold mt-2">Изменение категории</h1>
      </header>


      <!-- Radio Group for Categories -->
      <RadioGroup v-model="selected" class="flex-1 overflow-y-auto px-4 py-4 space-y-4">
        <RadioGroupOption v-for="category in store.categories" :key="category.id" :value="category.id"
          v-slot="{ checked }" class="flex items-center justify-between cursor-pointer px-2 py-3">
          <div class="flex items-center space-x-5">
            <span class="text-base">{{ category.name }}</span>
          </div>
          <div class="w-6 h-6 flex items-center justify-center">
            <CheckCircleIcon v-if="checked" class="w-6 h-6 text-blue-400" />
            <div v-else class="w-5 h-5 rounded-full border border-black/30"></div>
          </div>
        </RadioGroupOption>
      </RadioGroup>


      <!-- Footer -->
      <footer class="sticky bottom-0 bg-gray-100 p-8 rounded-t-2xl">
        <button @click="changeCategory"
          class="w-full bg-yellow-400 text-black py-4 rounded-2xl action:bg-yellow-300 transition">
          Поменять категорию
        </button>
      </footer>
    </div>
  </template>
</template>
