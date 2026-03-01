-- +goose Up
CREATE TABLE race_stages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    stage_number INT NOT NULL,
    stage_distance NUMERIC NOT NULL,
    stage_time INT NOT NULL,
    stage_position INT NOT NULL,
    note TEXT,
    race_id uuid NOT NULL,
    FOREIGN KEY (race_id) REFERENCES race_results(id) ON DELETE CASCADE
    
);
CREATE INDEX idx_race_stages_race_id ON race_stages(race_id);
-- +goose Down
DROP TABLE IF EXISTS race_stages;