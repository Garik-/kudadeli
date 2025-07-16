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
)

func TestExpenseCRUD(t *testing.T) {
	ctx := context.Background()

	tmpFile := "test_expenses.db"
	defer os.Remove(tmpFile)

	srv, err := database.New(ctx, tmpFile)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
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
		if err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	})

	t.Run("List and Check Insert", func(t *testing.T) {
		items, err := srv.List(ctx)
		if err != nil {
			t.Fatalf("list failed: %v", err)
		}

		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(items))
		}

		got := items[0]
		if got.ID != expense.ID {
			t.Errorf("ID mismatch")
		}

		if !got.Amount.Equal(expense.Amount) {
			t.Errorf("amount mismatch")
		}
	})

	t.Run("Update", func(t *testing.T) {
		expense.Description = "Updated description"
		expense.Amount = decimal.NewFromFloat(222.22)
		expense.UpdatedAt = expense.UpdatedAt.Add(1 * time.Hour)

		err := srv.Update(ctx, expense)
		if err != nil {
			t.Fatalf("update failed: %v", err)
		}
	})

	t.Run("List and Check Update", func(t *testing.T) {
		items, err := srv.List(ctx)
		if err != nil {
			t.Fatalf("list after update failed: %v", err)
		}

		if items[0].Description != "Updated description" {
			t.Errorf("description not updated")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if err := srv.Delete(ctx, expense.ID); err != nil {
			t.Fatalf("delete failed: %v", err)
		}

		items, err := srv.List(ctx)
		if err != nil {
			t.Fatalf("list after delete failed: %v", err)
		}

		if len(items) != 0 {
			t.Errorf("expected 0 items after delete, got %d", len(items))
		}
	})
}
