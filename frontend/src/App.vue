<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useCategoriesStore } from '@/stories/categoriesStore'
import { RouterView } from 'vue-router'
import { useTmaStore } from './stories/tmaStore'

const categoriesStore = useCategoriesStore()
const tmaStore = useTmaStore()
const loading = computed(() => categoriesStore.loading)

const initialLoader = document.getElementById('initial-loader')
if (initialLoader) {
  watch(loading, async (value) => {
    if (value) {
      initialLoader.remove()
    }
  })
}

onMounted(() => {
  tmaStore.init()
  categoriesStore.loadCategories()
})
</script>

<template>
  <RouterView />
</template>
