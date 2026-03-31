CREATE TABLE IF NOT EXISTS tmux_launchers (
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL COLLATE NOCASE UNIQUE,
    icon         TEXT NOT NULL DEFAULT 'terminal',
    command      TEXT NOT NULL DEFAULT '',
    cwd_mode     TEXT NOT NULL DEFAULT 'session',
    cwd_value    TEXT NOT NULL DEFAULT '',
    window_name  TEXT NOT NULL DEFAULT '',
    sort_order   INTEGER NOT NULL DEFAULT 0,
    created_at   TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at   TEXT NOT NULL DEFAULT (datetime('now')),
    last_used_at TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_tmux_launchers_sort_order
    ON tmux_launchers (sort_order ASC, name ASC);
