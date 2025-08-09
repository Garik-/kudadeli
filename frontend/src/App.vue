<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useCategoriesStore } from '@/stories/categoriesStore';
import { RouterView } from 'vue-router';

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


</script>

<template>
  <RouterView />
</template>
