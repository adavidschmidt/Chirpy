-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;


-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1;


-- name: UpdateUser :one
UPDATE users
SET updated_at = NOW(), email = $1, hashed_password = $2
WHERE id = $3
RETURNING *;

-- name: AddChirpyRed :exec
UPDATE users
SET is_chirpy_red = 'True'
WHERE id = $1;

-- name: RemoveChirpyRed :exec
UPDATE users
SET is_chirpy_red = 'False'
WHERE id = $1;