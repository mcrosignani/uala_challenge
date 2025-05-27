CREATE TABLE IF NOT EXISTS followers (
    user_id INTEGER NOT NULL,
    follow_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, follow_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follow_id) REFERENCES users(id) ON DELETE CASCADE
);