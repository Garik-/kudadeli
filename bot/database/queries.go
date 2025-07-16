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
    payment_type_id INTEGER NOT NULL
)
`

	insertExpense = `
INSERT INTO expenses (
	id, created_at, updated_at, category_id, description, amount, payment_type_id
) VALUES (?, ?, ?, ?, ?, ?, ?)
`

	updateExpense = `
UPDATE expenses
SET updated_at = ?, category_id = ?, description = ?, amount = ?, payment_type_id = ?
WHERE id = ?
`

	deleteExpense = `DELETE FROM expenses WHERE id = ?`

	selectExpenses = `
SELECT id, created_at, updated_at, category_id, description, amount, payment_type_id
FROM expenses
ORDER BY created_at DESC
	`
)
