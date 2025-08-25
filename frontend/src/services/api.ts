import type { Expense } from '@/models/expense'
import type { Category } from '@/models/category'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

export async function fetchExpenses(): Promise<Expense[]> {
  const response = await fetch(`${API_BASE_URL}/v1/expenses`)
  return response.json()
}

export async function fetchCategories(): Promise<Category[]> {
  const response = await fetch(`${API_BASE_URL}/v1/categories`)
  return response.json()
}

function authorizationHeader(token?: string) {
  if (token) {
    return {
      Authorization: `tma ${token}`,
    }
  }
}

export async function updateExpenseCategory(expenseId: string, category: number, token?: string) {
  const response = await fetch(`${API_BASE_URL}/v1/expenses/${expenseId}/category`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      ...authorizationHeader(token),
    },
    body: JSON.stringify({ category }),
  })

  if (response.status === 204 || response.status === 403) {
    return response.status
  }

  throw new Error(`Unexpected response status: ${response.status}`)
}
