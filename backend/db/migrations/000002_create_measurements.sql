-- +goose Up
CREATE TABLE measurements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    date TIMESTAMPTZ NOT NULL,
    weight NUMERIC, 
    neck NUMERIC,
    chest NUMERIC,
    waist NUMERIC,
    hips NUMERIC,
    right_forearm NUMERIC,
    left_forearm NUMERIC,
    right_bicep NUMERIC,
    left_bicep NUMERIC,
    right_quad NUMERIC,
    left_quad NUMERIC,
    right_calf NUMERIC,
    left_calf NUMERIC,
    bodyfat_percentage NUMERIC,
    user_id uuid NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_measurements_user_id ON measurements(user_id);
-- +goose Down
DROP TABLE IF EXISTS measurements;