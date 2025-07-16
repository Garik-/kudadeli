package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ExpenseID = uuid.UUID

type Expense struct {
	ID          ExpenseID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Category    Category
	PaymentType PaymentType
	Description string
	Amount      decimal.Decimal
}
