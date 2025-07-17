<template>
  <div class="max-w-3xl mx-auto px-4 py-8 space-y-8">
    <!-- Filters -->
    <div class="flex flex-wrap gap-2">
      <Listbox v-model="selectedMonth">
        <div class="relative">
          <ListboxButton class="px-4 py-2 bg-blue-600 text-white rounded-md">{{ selectedMonth }}</ListboxButton>
          <ListboxOptions class="absolute mt-1 bg-white border rounded-md w-32 z-10">
            <ListboxOption v-for="month in months" :key="month" :value="month" class="px-4 py-2 hover:bg-gray-100 cursor-pointer">
              {{ month }}
            </ListboxOption>
          </ListboxOptions>
        </div>
      </Listbox>
      <button class="px-4 py-2 rounded-md bg-blue-500 text-white">Black</button>
      <button class="px-4 py-2 rounded-md bg-gray-100 text-gray-700">Без переводов</button>
    </div>

    <!-- Time Period Tabs -->
    <div class="flex gap-4 text-sm">
      <button class="font-medium text-gray-800">Неделя</button>
      <button class="text-gray-500">Месяц</button>
      <button class="text-gray-500">Год</button>
    </div>

    <!-- Category Summary -->
    <div class="flex flex-wrap gap-2">
      <div v-for="category in categories" :key="category.name" class="flex items-center gap-2 px-3 py-1 rounded-full text-sm font-medium" :class="category.bg">
        <component :is="category.icon" class="w-4 h-4" />
        {{ category.name }} {{ category.amount }} ₽
      </div>
    </div>

    <!-- Transactions by Date -->
    <div v-for="group in groupedTransactions" :key="group.date" class="space-y-4">
      <div class="text-lg font-bold">{{ group.date }}</div>
      <div class="space-y-4">
        <div v-for="tx in group.items" :key="tx.id" class="flex justify-between items-center">
          <div class="flex items-center gap-4">
            <img :src="tx.logo" alt="" class="w-10 h-10 rounded-full" />
            <div>
              <div class="font-medium">{{ tx.title }}</div>
              <div class="text-sm text-gray-500">{{ tx.subtitle }}</div>
            </div>
          </div>
          <div class="text-right">
            <div class="text-red-600 font-medium">-{{ tx.amount }} ₽</div>
            <div class="text-sm text-gray-500">{{ tx.card }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Listbox, ListboxButton, ListboxOptions, ListboxOption } from '@headlessui/vue'
import { BanknotesIcon, ShoppingBagIcon, FireIcon, TruckIcon, EllipsisHorizontalIcon } from '@heroicons/vue/24/solid'

const selectedMonth = ref('Июль')
const months = ['Июль', 'Июнь', 'Май']

const categories = [
  { name: 'Переводы', amount: 332117, icon: BanknotesIcon, bg: 'bg-blue-100 text-blue-800' },
  { name: 'Маркетплейсы', amount: 4906, icon: ShoppingBagIcon, bg: 'bg-pink-100 text-pink-800' },
  { name: 'Рестораны', amount: 4323, icon: FireIcon, bg: 'bg-red-100 text-red-800' },
  { name: 'Фастфуд', amount: 3790, icon: FireIcon, bg: 'bg-yellow-100 text-yellow-800' },
  { name: 'Такси', amount: 2813, icon: TruckIcon, bg: 'bg-yellow-200 text-yellow-900' },
  { name: 'Остальное', amount: 4363, icon: EllipsisHorizontalIcon, bg: 'bg-gray-100 text-gray-800' },
]

const groupedTransactions = [
  {
    date: 'Вчера',
    items: [
      { id: 1, title: 'Пятёрочка', subtitle: 'Супермаркеты', amount: '1582,9', card: 'Дебетовая карта', logo: '/logos/pyaterochka.png' },
      { id: 2, title: 'Пятёрочка', subtitle: 'Супермаркеты', amount: '59,99', card: 'Дебетовая карта', logo: '/logos/pyaterochka.png' },
      { id: 3, title: 'Павел К.', subtitle: 'Переводы', amount: '120000', card: 'Black', logo: '/logos/person-yellow.png' },
      { id: 4, title: 'Анфиса Р.', subtitle: 'Переводы', amount: '50000', card: 'Black', logo: '/logos/person-yellow.png' },
    ],
  },
  {
    date: '16 июля',
    items: [
      { id: 5, title: 'Яндекс Такси', subtitle: 'Такси', amount: '220', card: 'Дебетовая карта', logo: '/logos/yandex-taxi.png' },
      { id: 6, title: 'Яндекс Такси', subtitle: 'Такси', amount: '1389', card: 'Дебетовая карта', logo: '/logos/yandex-taxi.png' },
    ],
  },
  {
    date: '15 июля',
    items: [
      { id: 7, title: 'Яндекс Такси', subtitle: 'Такси', amount: '1204', card: '', logo: '/logos/yandex-taxi.png' },
    ],
  },
]
</script>

<style scoped>
</style>
