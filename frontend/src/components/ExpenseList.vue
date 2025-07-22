<template>
  <div class="max-w-3xl mx-auto px-4 py-8 space-y-8">

     <div className="flex gap-6 p-6 bg-gray-50 rounded-2xl  min-h-[120px]">

      <div className="flex flex-col justify-between bg-white p-6 rounded-2xl shadow-lg w-1/2">
        <div className="text-2xl font-extrabold leading-none">{{ totalAmount }}</div>
        <div className="text-gray-500 mb-3">Траты</div>
        <!--<div className="flex h-4 w-full rounded-full overflow-hidden">

          <div className="bg-blue-400 flex-grow"></div>

          <div className="bg-indigo-300 w-6"></div>
          <div className="bg-rose-500 w-6"></div>
          <div className="bg-pink-500 w-6"></div>
          <div className="bg-amber-400 w-6"></div>
          <div className="bg-slate-400 w-6"></div>
        </div> -->
      </div>

      <div className="flex flex-col justify-between bg-white p-6 rounded-2xl shadow-lg w-1/2">
        <div className="text-2xl font-extrabold leading-none">{{  budgetAmount }}</div>
        <div className="text-gray-500 mb-3">Бюджет</div>
        <!--<div className="flex h-4 w-full rounded-full overflow-hidden">
          <div className="bg-blue-400 flex-grow"></div>
          <div className="bg-sky-600 w-6"></div>
          <div className="bg-teal-300 w-6"></div>
        </div>-->
      </div>
    </div>



    <!-- Transactions by Date -->
    <div  v-for="group in groupedTransactions" :key="group.date" class="space-y-4">
      <div class="text-lg font-bold">{{ group.date }}</div>
      <div class="space-y-4">
        <div v-for="tx in group.items" :key="tx.id" class="flex justify-between items-center">
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
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Listbox, ListboxButton, ListboxOptions, ListboxOption } from '@headlessui/vue'
import { BanknotesIcon, ShoppingBagIcon, FireIcon, TruckIcon, EllipsisHorizontalIcon } from '@heroicons/vue/24/solid'

import { isToday, isYesterday } from 'date-fns'

interface ExpenseResponse {
  id: string;
  createdAt: string;
  updatedAt: string;
  category: string;
  paymentType: string;
  description: string;
  amount: string;
  userId: number;
}

const BUDGET = 3_000_000.00

const selectedMonth = ref('Июль')
const months = ['Июль', 'Июнь', 'Май']


const totalAmount = ref('')
const budgetAmount = ref('')

function getTotalAmount(data: ExpenseResponse[]) {
  let amount = 0;

   data.forEach(expense => {
      amount += parseFloat(expense.amount)
   })

  return amount
}

const categories = [
  { name: 'Переводы', amount: 332117, icon: BanknotesIcon, bg: 'bg-blue-100 text-blue-800' },
  { name: 'Маркетплейсы', amount: 4906, icon: ShoppingBagIcon, bg: 'bg-pink-100 text-pink-800' },
  { name: 'Рестораны', amount: 4323, icon: FireIcon, bg: 'bg-red-100 text-red-800' },
  { name: 'Фастфуд', amount: 3790, icon: FireIcon, bg: 'bg-yellow-100 text-yellow-800' },
  { name: 'Такси', amount: 2813, icon: TruckIcon, bg: 'bg-yellow-200 text-yellow-900' },
  { name: 'Остальное', amount: 4363, icon: EllipsisHorizontalIcon, bg: 'bg-gray-100 text-gray-800' },
]

const groupedTransactions = ref([])

function formatDateToGroupLabel(dateString) {
  const date = new Date(dateString)


  if (isToday(date)) return 'Сегодня'
  if (isYesterday(date)) return 'Вчера'

  return date.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long'
  }) // например: 16 июля
}

function capitalizeFirstLetter(str: string) {
  if (typeof str !== 'string' || str.length === 0) return str;

  // Удаляем ведущие пробелы, но запоминаем их для восстановления
  const leadingSpaces = str.match(/^\s*/)[0];
  const trimmed = str.trimStart();

  if (trimmed.length === 0) return str; // строка из одних пробелов

  // Берём первый символ с учётом суррогатных пар (UTF-16)
  const firstChar = [...trimmed][0];
  const firstCharUpper = firstChar.toLocaleUpperCase();

  // Остальная часть строки после первого символа
  const rest = [...trimmed].slice(1).join('');

  return leadingSpaces + firstCharUpper + rest;
}

function formatPrice(amount: number): string {
  return amount.toLocaleString('ru-RU', { style: 'currency', currency: 'RUB', minimumFractionDigits: 0,maximumFractionDigits: 2})
}


function transformExpenses(data: ExpenseResponse[]) {
  const map = new Map()

  data.forEach(expense => {
    const groupKey = formatDateToGroupLabel(expense.createdAt)

    if (!map.has(groupKey)) {
      map.set(groupKey, [])
    }

    map.get(groupKey).push({
      id: expense.id,
      title: capitalizeFirstLetter(expense.description),
      category: expense.category,
      amount: formatPrice(parseFloat(expense.amount)),
      paymentType: expense.paymentType,
      // logo: getLogo(expense), // логика для логотипа
    })
  })

  // Преобразуем Map в массив
  return Array.from(map.entries()).map(([date, items]) => ({
    date,
    items,
  }))
}

onMounted(async () => {
  try {
    const res = await fetch('http://localhost:8080/v1/expenses')
    if (!res.ok) throw new Error('Ошибка загрузки')

    const data = await res.json()
    groupedTransactions.value = transformExpenses(data)

    const amount = getTotalAmount(data)

    totalAmount.value = formatPrice(amount)
    budgetAmount.value = formatPrice(BUDGET - amount)
  } catch (err) {
    console.error(err)
   // error.value = err.message
  }
})

/*
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
]*/
</script>

<style scoped>
</style>
