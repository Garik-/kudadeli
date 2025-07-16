package parser

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"kudadeli/model"
)

var (
	ErrEmptyMessage        = errors.New("empty message")
	ErrNotEnoughData       = errors.New("not enough data")
	ErrPaymentTypeNotFound = errors.New("payment type not found")
	ErrAmountNotFound      = errors.New("amount not found")

	categoryWords = map[string]model.Category{ //nolint:gochecknoglobals
		"материалы":   model.CategoryMaterials,
		"услуги":      model.CategoryLabor,
		"инструменты": model.CategoryTools,
		"мебель":      model.CategoryFurniture,
	}
)

const minWords = 2

func Message(input string) (model.Expense, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return model.Expense{}, ErrEmptyMessage
	}

	words := strings.Fields(input)
	if len(words) < minWords {
		return model.Expense{}, ErrNotEnoughData
	}

	var (
		paymentType      model.PaymentType
		category         = model.CategoryUnexpected
		amount           decimal.Decimal
		foundPaymentType bool
		foundAmount      bool
		descriptionWords = make([]string, 0, len(words)-minWords)
	)

	for i := range words {
		word := words[i]

		// Платеж
		if !foundPaymentType {
			switch word {
			case "нал", "наличные":
				paymentType = model.PaymentTypeCash
				foundPaymentType = true

				continue
			case "карта":
				paymentType = model.PaymentTypeCard
				foundPaymentType = true

				continue
			}
		}

		// Сумма (после типа оплаты)
		if foundPaymentType && !foundAmount {
			if amt, err := decimal.NewFromString(word); err == nil {
				amount = amt
				foundAmount = true

				continue
			}
		}

		// Категория
		if cat, ok := categoryWords[word]; ok {
			category = cat

			continue
		}

		// Все остальное — описание
		descriptionWords = append(descriptionWords, word)
	}

	if !foundPaymentType {
		return model.Expense{}, ErrPaymentTypeNotFound
	}

	if !foundAmount {
		return model.Expense{}, ErrAmountNotFound
	}

	createdAt := time.Now()

	return model.Expense{
		ID:          uuid.New(),
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
		Category:    category,
		PaymentType: paymentType,
		Description: strings.Join(descriptionWords, " "),
		Amount:      amount,
	}, nil
}
