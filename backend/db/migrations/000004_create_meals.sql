-- +goose Up
CREATE TABLE meals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    date TIMESTAMPTZ NOT NULL,
    food_name TEXT,
    meal_type TEXT NOT NULL,
    quantity NUMERIC NOT NULL,
    protein NUMERIC,
    carbs NUMERIC,
    fat NUMERIC,
    estimated_calories INT,
    note TEXT,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    
);
CREATE INDEX idx_meals_user_id ON meals(user_id);
-- +goose Down
DROP TABLE IF EXISTS meals;