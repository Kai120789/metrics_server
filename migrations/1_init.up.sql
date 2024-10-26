CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    type VARCHAR(50),
    value FLOAT NULL,
    delta INTEGER NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
)