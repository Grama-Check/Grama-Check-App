// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: gramacheck.sql

package db

import (
	"context"
)

const createCheck = `-- name: CreateCheck :one
INSERT INTO checks (
    nic,
    address,
    name,
    email,
    idcheck,
    addresscheck,
    policecheck,
    failed
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8)
RETURNING nic, name, address, email, idcheck, addresscheck, policecheck, failed
`

type CreateCheckParams struct {
	Nic          string `json:"nic"`
	Address      string `json:"address"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Idcheck      bool   `json:"idcheck"`
	Addresscheck bool   `json:"addresscheck"`
	Policecheck  bool   `json:"policecheck"`
	Failed       bool   `json:"failed"`
}

func (q *Queries) CreateCheck(ctx context.Context, arg CreateCheckParams) (Check, error) {
	row := q.db.QueryRowContext(ctx, createCheck,
		arg.Nic,
		arg.Address,
		arg.Name,
		arg.Email,
		arg.Idcheck,
		arg.Addresscheck,
		arg.Policecheck,
		arg.Failed,
	)
	var i Check
	err := row.Scan(
		&i.Nic,
		&i.Name,
		&i.Address,
		&i.Email,
		&i.Idcheck,
		&i.Addresscheck,
		&i.Policecheck,
		&i.Failed,
	)
	return i, err
}

const deleteCheck = `-- name: DeleteCheck :exec
DELETE FROM checks 
WHERE nic = $1
`

func (q *Queries) DeleteCheck(ctx context.Context, nic string) error {
	_, err := q.db.ExecContext(ctx, deleteCheck, nic)
	return err
}

const getCheck = `-- name: GetCheck :one
SELECT nic, name, address, email, idcheck, addresscheck, policecheck, failed FROM checks
WHERE nic = $1
LIMIT 1
`

func (q *Queries) GetCheck(ctx context.Context, nic string) (Check, error) {
	row := q.db.QueryRowContext(ctx, getCheck, nic)
	var i Check
	err := row.Scan(
		&i.Nic,
		&i.Name,
		&i.Address,
		&i.Email,
		&i.Idcheck,
		&i.Addresscheck,
		&i.Policecheck,
		&i.Failed,
	)
	return i, err
}

const updateAddressCheck = `-- name: UpdateAddressCheck :exec
UPDATE checks SET addresscheck = true
WHERE nic = $1
`

func (q *Queries) UpdateAddressCheck(ctx context.Context, nic string) error {
	_, err := q.db.ExecContext(ctx, updateAddressCheck, nic)
	return err
}

const updateFailed = `-- name: UpdateFailed :exec
UPDATE checks SET failed = true
WHERE nic = $1
`

func (q *Queries) UpdateFailed(ctx context.Context, nic string) error {
	_, err := q.db.ExecContext(ctx, updateFailed, nic)
	return err
}

const updateIdentityCheck = `-- name: UpdateIdentityCheck :exec
UPDATE checks SET idcheck = true
WHERE nic = $1
`

func (q *Queries) UpdateIdentityCheck(ctx context.Context, nic string) error {
	_, err := q.db.ExecContext(ctx, updateIdentityCheck, nic)
	return err
}

const updatePoliceCheck = `-- name: UpdatePoliceCheck :exec
UPDATE checks SET policecheck = true
WHERE nic = $1
`

func (q *Queries) UpdatePoliceCheck(ctx context.Context, nic string) error {
	_, err := q.db.ExecContext(ctx, updatePoliceCheck, nic)
	return err
}
