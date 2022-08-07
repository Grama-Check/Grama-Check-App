-- name: GetUser :one
SELECT * FROM users
WHERE nic = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    nic,
    address,
    name,
    email,
    idcheck,
    addresscheck,
    policecheck,
    failed

) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8 )
RETURNING *;

-- name: UpdateID :exec
UPDATE users SET idcheck = true
WHERE nic = $1;

-- name: UpdateAddress :exec
UPDATE users SET addresscheck = true
WHERE nic = $1;

-- name: UpdatePolice :exec
UPDATE users SET policecheck = true
WHERE nic = $1;

-- name: UpdateFailed :exec
UPDATE users SET failed = $2
WHERE nic = $1;