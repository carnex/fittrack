-- name: CreateUser :one
INSERT INTO users ( username, email, password, reset_method, security_question, security_question_answer)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)RETURNING id, created_at, username, email, reset_method, use_metric;

-- name: GetUserByUsername :one
select id, username, password
from users
where username = $1;

-- name: GetUserByEmail :one
select id, email
from users
where email = $1;
 