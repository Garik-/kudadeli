package database

const (
	createExpenses = `
CREATE TABLE expenses (
    id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
    category_id INTEGER NOT NULL,
    description TEXT,
    amount TEXT NOT NULL,
    payment_type_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	deleted_at DATETIME
)
`

	insertExpense = `
INSERT INTO expenses (
	id, created_at, updated_at, category_id, description, amount, payment_type_id, user_id
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

	updateExpense = `
UPDATE expenses
SET updated_at = ?, category_id = ?, description = ?, amount = ?, payment_type_id = ?
WHERE id = ?
`

	deleteExpense = `UPDATE expenses SET deleted_at = datetime('now') WHERE id = ?`

	selectExpenses = `
SELECT id, created_at, updated_at, category_id, description, amount, payment_type_id, user_id
FROM expenses WHERE deleted_at IS NULL
ORDER BY created_at DESC
	`
)
