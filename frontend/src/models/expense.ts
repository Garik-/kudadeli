export interface Expense {
  id: string
  createdAt: string
  updatedAt: string
  category: string
  paymentType: string
  description: string
  amount: string
  userId: number
}

export function getTotalAmount(data: Expense[]) {
  let amount = 0

  data.forEach((expense) => {
    amount += parseFloat(expense.amount)
  })

  return amount
}
