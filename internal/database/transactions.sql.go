// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: transactions.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createUserTransaction = `-- name: CreateUserTransaction :one
INSERT INTO transactions (user_id, wallet_id, amount, type, description, transaction_date)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, wallet_id, amount, type, description, transaction_date, created_at, updated_at
`

type CreateUserTransactionParams struct {
	UserID          uuid.UUID      `json:"user_id"`
	WalletID        uuid.UUID      `json:"wallet_id"`
	Amount          pgtype.Numeric `json:"amount"`
	Type            string         `json:"type"`
	Description     string         `json:"description"`
	TransactionDate time.Time      `json:"transaction_date"`
}

func (q *Queries) CreateUserTransaction(ctx context.Context, arg CreateUserTransactionParams) (Transaction, error) {
	row := q.db.QueryRow(ctx, createUserTransaction,
		arg.UserID,
		arg.WalletID,
		arg.Amount,
		arg.Type,
		arg.Description,
		arg.TransactionDate,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.WalletID,
		&i.Amount,
		&i.Type,
		&i.Description,
		&i.TransactionDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getManyUserTransaction = `-- name: GetManyUserTransaction :one
SELECT id, user_id, wallet_id, amount, type, description, transaction_date, created_at, updated_at FROM transactions WHERE user_id = $1
`

func (q *Queries) GetManyUserTransaction(ctx context.Context, userID uuid.UUID) (Transaction, error) {
	row := q.db.QueryRow(ctx, getManyUserTransaction, userID)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.WalletID,
		&i.Amount,
		&i.Type,
		&i.Description,
		&i.TransactionDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
