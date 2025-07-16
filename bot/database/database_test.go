package database_test

import (
	"context"
	"kudadeli/database"
	"kudadeli/model"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpenseCRUD(t *testing.T) {
	ctx := context.Background()

	tmpFile := "test_expenses.db"
	defer os.Remove(tmpFile)

	srv, err := database.New(ctx, tmpFile)
	require.NoError(t, err, "failed to create database")

	defer srv.Close()

	expense := model.Expense{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
		UpdatedAt:   time.Now().UTC().Truncate(time.Second),
		Category:    model.CategoryMaterials,
		PaymentType: model.PaymentTypeCard,
		Description: "Test expense",
		Amount:      decimal.NewFromFloat(123.45),
	}

	t.Run("Insert", func(t *testing.T) {
		err := srv.Insert(ctx, expense)
		require.NoError(t, err, "insert failed")
	})

	t.Run("List and Check Insert", func(t *testing.T) {
		items, err := srv.List(ctx)
		require.NoError(t, err, "list failed")

		require.Len(t, items, 1, "expected 1 item")

		got := items[0]
		assert.Equal(t, expense.ID, got.ID, "ID mismatch")
		assert.True(t, got.Amount.Equal(expense.Amount), "amount mismatch")
	})

	t.Run("Update", func(t *testing.T) {
		expense.Description = "Updated description"
		expense.Amount = decimal.NewFromFloat(222.22)
		expense.UpdatedAt = expense.UpdatedAt.Add(1 * time.Hour)

		err := srv.Update(ctx, expense)
		require.NoError(t, err, "update failed")
	})

	t.Run("List and Check Update", func(t *testing.T) {
		items, err := srv.List(ctx)
		require.NoError(t, err, "list after update failed")

		require.NotEmpty(t, items, "expected items after update")
		assert.Equal(t, "Updated description", items[0].Description, "description not updated")
	})

	t.Run("Delete", func(t *testing.T) {
		err := srv.Delete(ctx, expense.ID)
		require.NoError(t, err, "delete failed")

		items, err := srv.List(ctx)
		require.NoError(t, err, "list after delete failed")

		assert.Empty(t, items, "expected 0 items after delete")
	})
}
