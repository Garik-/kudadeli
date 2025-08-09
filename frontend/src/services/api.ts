import axios from 'axios'
import type { Expense } from '@/models/expense'
import type { Category } from '@/models/category'
import { retrieveLaunchParams } from '@telegram-apps/sdk'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

export async function fetchExpenses() {
  const response = await axios.get<Expense[]>(`${API_BASE_URL}/v1/expenses`)
  return response.data
}

export async function fetchCategories() {
  const response = await axios.get<Category[]>(`${API_BASE_URL}/v1/categories`)
  return response.data
}

function authorizationHeader() {
  try {
    const { initDataRaw } = retrieveLaunchParams()
    return {
      Authorization: `tma ${initDataRaw}`,
    }
  } catch (e: unknown) {
    if (e instanceof Error) {
      console.error(e.message)
    }
  }

  return {}
}

export async function updateExpenseCategory(expenseId: string, category: number) {
  const response = await axios.put(
    `${API_BASE_URL}/v1/expenses/${expenseId}/category`,
    {
      category,
    },
    {
      validateStatus: (status) => status === 204 || status === 403,
      headers: authorizationHeader(),
    },
  )

  return response.status
}
