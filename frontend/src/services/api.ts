import axios from 'axios'
import type { Expense } from '@/models/expense'
import type { Category } from '@/models/category'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

export async function fetchExpenses() {
  const response = await axios.get<Expense[]>(`${API_BASE_URL}/v1/expenses`)
  return response.data
}

/*
  const res = await fetch(`${baseUrl}/v1/expenses`)
  if (!res.ok) throw new LoadingError()

  const data = await res.json()
  groupedTransactions.value = transformExpenses(data)
  groupedAmount.value = transformExpensesAmount(data)

  const amount = getTotalAmount(data)

  totalAmount.value = formatPrice(amount)
  budgetAmount.value = formatPrice(BUDGET - amount)
}
  */

export async function fetchCategories() {
  const response = await axios.get<Category[]>(`${API_BASE_URL}/v1/categories`)
  return response.data
}
