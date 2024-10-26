CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    type VARCHAR(50),
    value FLOAT,
    delta INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW()
)