-- +goose Up
CREATE TABLE integration_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider TEXT NOT NULL,
    provider_user_id TEXT
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    
);
CREATE INDEX idx_integration_tokens_user_id ON integration_tokens(user_id);
CREATE UNIQUE INDEX idx_integration_tokens_user_provider 
    ON integration_tokens(user_id, provider);
-- +goose Down
DROP TABLE IF EXISTS integration_tokens;