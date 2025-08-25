<script setup lang="ts">
import FilterButton from './FilterButton.vue'
import { useFiltersStore } from '@/stories/filtersStore'
import { capitalizeFirstLetter } from '@/utils/formatter'
import { ref, computed } from 'vue'
import type { Filter } from '@/stories/filtersStore'
const filtersStore = useFiltersStore()
const isCart = ref(false)
const isCategory = ref(true)

function handleFilterCard() {
  if (isCart.value) {
    filtersStore.setFilter('paymentType', 'карта')
  } else {
    filtersStore.removeFilter('paymentType')
  }
}

const categoryButton = computed(() => createButton(filtersStore.filter))

function createButton(filter: Filter) {
  if (!('category' in filter)) {
    return undefined
  }
  return {
    label: capitalizeFirstLetter(filter.category!),
  }
}

function removeCategoryFilter() {
  filtersStore.removeFilter('category')
  isCategory.value = true
}
</script>
<template>
  <div class="flex flex-wrap gap-2 mb-6">
    <FilterButton label="Картой" v-model="isCart" @change="handleFilterCard" />
    <FilterButton
      v-if="categoryButton"
      :label="categoryButton.label"
      v-model="isCategory"
      @click="removeCategoryFilter"
    />
  </div>
</template>
