CREATE TABLE IF NOT EXISTS history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    last_used INTEGER NOT NULL,
    use_count INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_history_name ON history(name);
CREATE INDEX IF NOT EXISTS idx_history_last_used ON history(last_used);
