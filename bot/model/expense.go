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
