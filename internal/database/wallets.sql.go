// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: wallets.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createUserWallet = `-- name: CreateUserWallet :one
INSERT INTO wallets (user_id, name, balance)
VALUES ($1, $2, $3)
RETURNING id, user_id, name, balance, currency, created_at, updated_at
`

type CreateUserWalletParams struct {
	UserID  uuid.UUID      `json:"user_id"`
	Name    string         `json:"name"`
	Balance pgtype.Numeric `json:"balance"`
}

func (q *Queries) CreateUserWallet(ctx context.Context, arg CreateUserWalletParams) (Wallet, error) {
	row := q.db.QueryRow(ctx, createUserWallet, arg.UserID, arg.Name, arg.Balance)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUserWallet = `-- name: DeleteUserWallet :one
DELETE FROM wallets WHERE id = $1
RETURNING id, user_id, name, balance, currency, created_at, updated_at
`

func (q *Queries) DeleteUserWallet(ctx context.Context, id uuid.UUID) (Wallet, error) {
	row := q.db.QueryRow(ctx, deleteUserWallet, id)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const depositUserWallet = `-- name: DepositUserWallet :one
UPDATE wallets
SET balance = balance + $2
WHERE id = $1 AND user_id = $3
RETURNING id, user_id, name, balance, currency, created_at, updated_at
`

type DepositUserWalletParams struct {
	ID      uuid.UUID      `json:"id"`
	Balance pgtype.Numeric `json:"balance"`
	UserID  uuid.UUID      `json:"user_id"`
}

func (q *Queries) DepositUserWallet(ctx context.Context, arg DepositUserWalletParams) (Wallet, error) {
	row := q.db.QueryRow(ctx, depositUserWallet, arg.ID, arg.Balance, arg.UserID)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getManyUserWallet = `-- name: GetManyUserWallet :many
SELECT id, user_id, name, balance, currency, created_at, updated_at FROM wallets WHERE user_id = $1
`

func (q *Queries) GetManyUserWallet(ctx context.Context, userID uuid.UUID) ([]Wallet, error) {
	rows, err := q.db.Query(ctx, getManyUserWallet, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Wallet
	for rows.Next() {
		var i Wallet
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserWalletById = `-- name: GetUserWalletById :one
SELECT id, user_id, name, balance, currency, created_at, updated_at FROM wallets WHERE id = $1
`

func (q *Queries) GetUserWalletById(ctx context.Context, id uuid.UUID) (Wallet, error) {
	row := q.db.QueryRow(ctx, getUserWalletById, id)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
