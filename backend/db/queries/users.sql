-- name: CreateUser :one
INSERT INTO users ( username, email, first, last, password, reset_method, security_question, security_question_answer)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)RETURNING *;

-- name: GetUserByUsername :one
select id, username, password
from users
where username = $1;

 