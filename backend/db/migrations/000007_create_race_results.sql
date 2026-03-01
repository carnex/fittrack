-- +goose Up
CREATE TABLE race_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    race_type TEXT NOT NULL,
    age_category TEXT, 
    age_category_position TEXT,
    total_finishers INT NOT NULL,
    total_finishers_category INT,
    type_multistage BOOLEAN DEFAULT FALSE,
    event_name TEXT NOT NULL,
    race_date TIMESTAMPTZ NOT NULL,
    distance DECIMAL NOT NULL,
    overall_position INT NOT NULL,
    finish_time INT NOT NULL,
    bib TEXT,
    notes TEXT,
    user_id uuid NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    
);
CREATE INDEX idx_race_results_user_id ON race_results(user_id);
-- +goose Down
DROP TABLE IF EXISTS race_results;