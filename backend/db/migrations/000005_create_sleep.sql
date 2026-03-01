-- +goose Up
CREATE TABLE sleep (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    start TIMESTAMPTZ NOT NULL,
    end TIMESTAMPTZ NOT NULL,
    quality INT,
    note TEXT,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    
);
CREATE INDEX idx_sleep_user_id ON sleep(user_id);
-- +goose Down
DROP TABLE IF EXISTS sleep;