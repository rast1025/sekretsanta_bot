CREATE TABLE IF NOT EXISTS users
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    chat_id   TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_username ON users (username);

CREATE TABLE IF NOT EXISTS groups
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    group_id   TEXT NOT NULL,
    UNIQUE (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX IF NOT EXISTS idx_group ON groups (group_id);