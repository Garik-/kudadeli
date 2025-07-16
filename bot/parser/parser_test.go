package parser_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"kudadeli/model"
	"kudadeli/parser"
)

func TestMessage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError error
		want        model.Expense
	}{
		{
			name:  "наличные, материалы, описание",
			input: "нал 1500 обои кухня",
			want: model.Expense{
				PaymentType: model.PaymentTypeCash,
				Category:    model.CategoryUnexpected,
				Amount:      decimal.NewFromInt(1500),
				Description: "обои кухня",
			},
		},
		{
			name:  "карта, двери, ванная",
			input: "карта 7000 двери ванная",
			want: model.Expense{
				PaymentType: model.PaymentTypeCard,
				Category:    model.CategoryUnexpected,
				Amount:      decimal.NewFromInt(7000),
				Description: "двери ванная",
			},
		},
		{
			name:  "карта, двери (одно слово в описании)",
			input: "карта 3200 двери",
			want: model.Expense{
				PaymentType: model.PaymentTypeCard,
				Category:    model.CategoryUnexpected,
				Amount:      decimal.NewFromInt(3200),
				Description: "двери",
			},
		},
		{
			name:  "услуги, демонтаж",
			input: "нал 5000 услуги демонтаж",
			want: model.Expense{
				PaymentType: model.PaymentTypeCash,
				Category:    model.CategoryLabor,
				Amount:      decimal.NewFromInt(5000),
				Description: "демонтаж",
			},
		},
		{
			name:        "пустая строка",
			input:       "",
			expectError: parser.ErrEmptyMessage,
		},
		{
			name:        "без суммы",
			input:       "нал краска",
			expectError: parser.ErrAmountNotFound,
		},
		{
			name:        "без типа оплаты",
			input:       "1200 обои",
			expectError: parser.ErrPaymentTypeNotFound,
		},
		{
			name:        "не число в поле суммы",
			input:       "карта тысяча обои",
			expectError: parser.ErrAmountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expense, err := parser.Message(tt.input)
			if tt.expectError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectError, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.PaymentType, expense.PaymentType)
			require.Equal(t, tt.want.Category, expense.Category)
			require.True(t, tt.want.Amount.Equal(expense.Amount))
			require.Equal(t, tt.want.Description, expense.Description)

			// Проверим что ID и время установлены
			require.NotEqual(t, uuid.Nil, expense.ID)
			require.WithinDuration(t, time.Now(), expense.CreatedAt, time.Second)
			require.WithinDuration(t, time.Now(), expense.UpdatedAt, time.Second)
		})
	}
}
