<script setup lang="ts">
import { defineProps } from 'vue';
import { XMarkIcon } from '@heroicons/vue/16/solid'
const props = defineProps({
  modelValue: { type: Boolean, default: false }, // для v-model
  label: { type: String, default: '' }
})

const emit = defineEmits(['update:modelValue', 'change'])

function toggle() {
  const newValue = !props.modelValue
  emit('update:modelValue', newValue)
  emit('change', newValue)
}

</script>
<template>
  <button @click="toggle" v-if="modelValue"
    class="filter-active inline-flex items-center justify-between gap-1 py-2 rounded-2xl bg-blue-500 text-white text-sm font-semibold cursor-pointer"><span>{{
      label }}</span>
    <XMarkIcon class="w-[1.5em] h-[1.5em]" />

  </button>
  <button @click="toggle" v-else
    class="filter py-2 rounded-2xl bg-gray-100 text-gray-700 text-sm font-bold cursor-pointer">{{ label
    }}</button>
</template>
<style>
.filter-active {
  padding-left: 12px;
  padding-right: 6px;
}

.filter {
  padding-left: 12px;
  padding-right: 12px;
}
</style>
