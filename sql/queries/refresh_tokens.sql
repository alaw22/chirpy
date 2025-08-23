-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 days', -- adds 60 days to the current time
    NULL
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.id FROM refresh_tokens
LEFT JOIN users
ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
    AND NOW() < refresh_tokens.expires_at
    AND refresh_tokens.revoked_at IS NULL;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1;
