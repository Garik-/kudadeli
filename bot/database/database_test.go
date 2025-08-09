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
		UserID:      1,
	}

	t.Run("Insert", func(t *testing.T) {
		err := srv.Insert(ctx, expense)
		require.NoError(t, err, "insert failed")
	})

	t.Run("List and Check Insert", func(t *testing.T) {
		items, err := srv.List(ctx, -1)
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
		items, err := srv.List(ctx, -1)
		require.NoError(t, err, "list after update failed")

		require.NotEmpty(t, items, "expected items after update")
		assert.Equal(t, "Updated description", items[0].Description, "description not updated")
	})

	t.Run("Delete", func(t *testing.T) {
		err := srv.Delete(ctx, expense.ID)
		require.NoError(t, err, "delete failed")

		items, err := srv.List(ctx, -1)
		require.NoError(t, err, "list after delete failed")

		assert.Empty(t, items, "expected 0 items after delete")
	})

	t.Run("UpdateCategory", func(t *testing.T) {
		// Insert a new expense to update its category
		expense := model.Expense{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC().Truncate(time.Second),
			UpdatedAt:   time.Now().UTC().Truncate(time.Second),
			Category:    model.CategoryMaterials,
			PaymentType: model.PaymentTypeCard,
			Description: "Category update test",
			Amount:      decimal.NewFromFloat(50.00),
			UserID:      2,
		}
		err := srv.Insert(ctx, expense)
		require.NoError(t, err, "insert for UpdateCategory failed")

		// Update the category
		newCategory := model.CategoryFurniture
		err = srv.UpdateCategory(ctx, expense.ID, newCategory)
		require.NoError(t, err, "updateCategory failed")

		// Verify the category was updated
		items, err := srv.List(ctx, -1)
		require.NoError(t, err, "list after UpdateCategory failed")
		require.NotEmpty(t, items, "expected items after UpdateCategory")
		assert.Equal(t, newCategory, items[0].Category, "category not updated")

		// Clean up
		err = srv.Delete(ctx, expense.ID)
		require.NoError(t, err, "delete after UpdateCategory failed")
	})

	t.Log("TODO: add LatestUpdatedAt")
}
