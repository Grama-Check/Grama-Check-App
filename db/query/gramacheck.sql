-- name: CreateCheck :one
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
RETURNING *;

-- name: GetCheck :one
SELECT * FROM checks
WHERE nic = $1
LIMIT 1;

-- name: UpdateIdentityCheck :exec
UPDATE checks SET idcheck = true
WHERE nic = $1;

-- name: UpdateAddressCheck :exec
UPDATE checks SET addresscheck = true
WHERE nic = $1;

-- name: UpdatePoliceCheck :exec
UPDATE checks SET policecheck = true
WHERE nic = $1;

-- name: UpdateFailed :exec
UPDATE checks SET failed = true
WHERE nic = $1;

-- name: DeleteCheck :exec
DELETE FROM checks 
WHERE nic = $1;