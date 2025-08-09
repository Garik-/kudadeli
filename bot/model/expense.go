package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ExpenseID = uuid.UUID

type Expense struct {
	ID          ExpenseID       `json:"id"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	Category    Category        `json:"category"`
	PaymentType PaymentType     `json:"paymentType"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	UserID      int64           `json:"userId"`
}

type Expenses []Expense

func (expenses Expenses) LatestUpdatedAt() time.Time {
	if len(expenses) == 0 {
		return time.Time{}
	}

	latest := expenses[0].UpdatedAt
	for i := 1; i < len(expenses); i++ {
		if expenses[i].UpdatedAt.After(latest) {
			latest = expenses[i].UpdatedAt
		}
	}

	return latest
}
