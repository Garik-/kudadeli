package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/shopspring/decimal"
	// Import SQLite driver anonymously for side-effects (registration with database/sql).
	_ "modernc.org/sqlite"

	"kudadeli/model"
)

type Service struct {
	db *sql.DB
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func New(ctx context.Context, database string) (*Service, error) {
	slog.InfoContext(ctx, "open database", "path", database)

	isExist := fileExists(database)
	slog.Debug("fileExists", "isExist", isExist)

	db, err := sql.Open("sqlite", database+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("sql open error: %w", err)
	}

	srv := &Service{db: db}

	if !isExist {
		err := srv.create(ctx)
		if err != nil {
			return nil, fmt.Errorf("create database: %w", err)
		}
	}

	return srv, nil
}

func (s *Service) Close() error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("database close error: %w", err)
	}

	return nil
}

func (s *Service) LatestUpdatedAt(ctx context.Context) (time.Time, error) {
	var updatedAt string

	err := s.db.QueryRowContext(ctx, latestUpdatedAt).Scan(&updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, nil // Нет записей
		}

		return time.Time{}, fmt.Errorf("failed to get latest updated_at: %w", err)
	}

	t, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	return t, nil
}

func (s *Service) Insert(ctx context.Context, expense model.Expense) error {
	_, err := s.db.ExecContext(ctx, insertExpense,
		expense.ID.String(),
		expense.CreatedAt.Format(time.RFC3339),
		expense.UpdatedAt.Format(time.RFC3339),
		int(expense.Category),
		expense.Description,
		expense.Amount.String(),
		int(expense.PaymentType),
		expense.UserID,
	)
	if err != nil {
		return fmt.Errorf("insert expense: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, expense model.Expense) error {
	_, err := s.db.ExecContext(ctx, updateExpense,
		expense.UpdatedAt.Format(time.RFC3339),
		int(expense.Category),
		expense.Description,
		expense.Amount.String(),
		int(expense.PaymentType),
		expense.ID.String(),
	)
	if err != nil {
		return fmt.Errorf("update expense: %w", err)
	}

	return nil
}

func (s *Service) UpdateCategory(ctx context.Context, expenseID model.ExpenseID, category model.Category) error {
	_, err := s.db.ExecContext(ctx, updateExpenseCategory,
		time.Now().Format(time.RFC3339),
		int(category),
		expenseID.String(),
	)
	if err != nil {
		return fmt.Errorf("update expense category: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id model.ExpenseID) error {
	_, err := s.db.ExecContext(ctx, deleteExpense, id.String())
	if err != nil {
		return fmt.Errorf("delete expense: %w", err)
	}

	return nil
}

func (s *Service) List(ctx context.Context, limit int) ([]model.Expense, error) {
	query := selectExpenses
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	slog.DebugContext(ctx, "list expenses", "limit", limit)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select expenses: %w", err)
	}
	defer rows.Close()

	var expenses []model.Expense

	for rows.Next() {
		var (
			expense                   model.Expense
			createdAt, updatedAt      string
			amountStr                 string
			categoryID, paymentTypeID int
			userID                    int64
		)

		err := rows.Scan(
			&expense.ID,
			&createdAt,
			&updatedAt,
			&categoryID,
			&expense.Description,
			&amountStr,
			&paymentTypeID,
			&userID,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		expense.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse created at: %w", err)
		}

		expense.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			return nil, fmt.Errorf("parse updated at: %w", err)
		}

		expense.Category = model.Category(categoryID)
		expense.PaymentType = model.PaymentType(paymentTypeID)

		expense.Amount, err = decimal.NewFromString(amountStr)
		if err != nil {
			return nil, fmt.Errorf("parse amount: %w", err)
		}

		expense.UserID = userID

		expenses = append(expenses, expense)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return expenses, nil
}

func (s *Service) create(ctx context.Context) error {
	slog.InfoContext(ctx, "create database")

	_, err := s.db.ExecContext(ctx, createExpenses)
	if err != nil {
		return fmt.Errorf("createExpenses: %w", err)
	}

	return nil
}
