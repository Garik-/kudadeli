
<script setup lang="ts">
import { ref, onMounted } from 'vue'

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

interface Item {
      id: string;
      title: string;
      category: string;
      amount:string;
      paymentType: string;
}

interface GroupedItems {
  date: string,
  items: Item[]
}

const BUDGET = 3_000_000.00


const totalAmount = ref('')
const budgetAmount = ref('')

function getTotalAmount(data: ExpenseResponse[]) {
  let amount = 0;

   data.forEach(expense => {
      amount += parseFloat(expense.amount)
   })

  return amount
}

const groupedTransactions = ref<GroupedItems[]>([])
const groupedAmount = ref<Record<string, string>>({})

function formatDateToGroupLabel(dateString: string) {
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
  const leadingSpaces = str.match(/^\s*/)?.[0] ?? '';
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

function transformExpensesAmount(data: ExpenseResponse[]): Record<string, string> {
  const amounts: Record<string, number> = {}

  data.forEach(expense => {
    const groupKey = formatDateToGroupLabel(expense.createdAt)
    const amount = parseFloat(expense.amount)

    if (!amounts[groupKey]) {
      amounts[groupKey] = 0
    }

    amounts[groupKey] += amount
  })

  // Преобразуем объект amounts в объект с отформатированными значениями
  const formattedAmounts: Record<string, string> = {}
  for (const key in amounts) {
    formattedAmounts[key] = formatPrice(amounts[key]); // форматируем каждое значение
  }
  return formattedAmounts
}


function transformExpenses(data: ExpenseResponse[]): GroupedItems[] {
  const map: Map<string, Item[]> = new Map()

  data.forEach(expense => {
    const groupKey = formatDateToGroupLabel(expense.createdAt)

    if (!map.has(groupKey)) {
      map.set(groupKey, [])
    }

    map.get(groupKey)!.push({
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

const baseUrl = import.meta.env.VITE_API_BASE_URL || ''

onMounted(async () => {
  try {
    const res = await fetch(`${baseUrl}/v1/expenses`)
    if (!res.ok) throw new Error('Ошибка загрузки')

    const data = await res.json()
    groupedTransactions.value = transformExpenses(data)
    groupedAmount.value = transformExpensesAmount(data)

    const amount = getTotalAmount(data)

    totalAmount.value = formatPrice(amount)
    budgetAmount.value = formatPrice(BUDGET - amount)
  } catch (err) {
    console.error(err)
   // error.value = err.message
  }
})
</script>
<style scoped>
.shadow-item {
  box-shadow: rgba(0, 0, 0, 0.12) 0 6px 34px 0;
}
</style>
<template>
  <div class="max-w-3xl mx-auto px-4 pb-8">

    <div className="sticky top-0 py-8">
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
        <div className="font-bold text-lg">{{  budgetAmount }}</div>
        <div className="text-gray-500 text-sm">Бюджет</div>
        <!--<div className="flex h-4 w-full rounded-full overflow-hidden">
          <div className="bg-blue-400 flex-grow"></div>
          <div className="bg-sky-600 w-6"></div>
          <div className="bg-teal-300 w-6"></div>
        </div>-->
      </div>
    </div>
</div>

<div className="space-y-8">

    <!-- Transactions by Date -->
    <div  v-for="group in groupedTransactions" :key="group.date" class="space-y-4">
      <div class="flex items-center justify-between">
        <div class="text-lg font-bold">{{ group.date }}</div>
        <div class="text-right text-gray-400">-{{ groupedAmount[group.date] }}</div>
    </div>

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
  </div>
</template>



