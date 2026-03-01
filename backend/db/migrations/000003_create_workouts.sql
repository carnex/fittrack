-- +goose Up
CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    source TEXT NOT NULL DEFAULT 'manual',
    workout_type TEXT NOT NULL,
    date TIMESTAMPTZ NOT NULL, 
    duration int NOT NULL,
    avg_heart_rate int,
    calories int,
    notes text,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    
);
CREATE INDEX idx_workouts_user_id ON workouts(user_id);
-- +goose Down
DROP TABLE IF EXISTS workouts;