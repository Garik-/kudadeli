<script setup lang="ts">
import ExpenseList from '@/components/ExpenseList.vue';
import CategorySelector from './components/CategorySelector.vue';
import { ref, computed, onMounted, watch } from 'vue'
import { useCategoriesStore } from '@/stories/categoriesStore';

const store = useCategoriesStore()
const loading = computed(() => store.loading)

const initialLoader = document.getElementById('initial-loader')
if (initialLoader) {
  watch(loading, async (value) => {
    if (value) {
      initialLoader.remove()
    }
  })
}

onMounted(() => {
  store.loadCategories()
})


const currentPath = ref(window.location.hash)

window.addEventListener('hashchange', () => {
  currentPath.value = window.location.hash
})

const currentView = computed(() => {
  if (!currentPath.value) {
    return ExpenseList
  }

  return CategorySelector
})

</script>

<template>
  <component :is="currentView" />
</template>
