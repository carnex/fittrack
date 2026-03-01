-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    first TEXT,
    last TEXT,
    date_of_birth date ,
    password TEXT NOT NULL,
    reset_method BOOLEAN NOT NULL,
    security_question TEXT,
    security_question_answer TEXT,
    use_metric BOOLEAN DEFAULT TRUE,
    height NUMERIC
);

-- +goose Down
DROP TABLE IF EXISTS users CASCADE;