-- 000006_marker-patterns.sql: Configurable marker detection patterns.
-- Moves hardcoded marker keywords to a database table with CRUD support.
-- The watchtower pane content scanner uses these patterns to detect
-- timeline markers in pane output.

CREATE TABLE IF NOT EXISTS marker_patterns (
    id         TEXT PRIMARY KEY,
    pattern    TEXT NOT NULL,
    severity   TEXT NOT NULL DEFAULT 'warn',
    label      TEXT NOT NULL DEFAULT '',
    enabled    INTEGER NOT NULL DEFAULT 1,
    priority   INTEGER NOT NULL DEFAULT 50,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_marker_patterns_priority
    ON marker_patterns (priority ASC, id ASC);

-- Seed default error-level patterns.
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.panic',              'panic',              'error', 'Panic detected',              1, 10);
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.fatal',              'fatal',              'error', 'Fatal error',                 1, 11);
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.error',              'error',              'error', 'Error detected',              1, 12);
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.oom',                'out of memory',      'error', 'Out of memory',               1, 13);
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.segfault',           'segfault',           'error', 'Segmentation fault',          1, 14);
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.segmentation-fault', 'segmentation fault', 'error', 'Segmentation fault',          1, 15);
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.connection-refused', 'connection refused', 'error', 'Connection refused',          1, 16);
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.killed',             'killed',             'error', 'Process killed',              1, 17);

-- Seed default warn-level patterns.
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.timeout',            'timeout',            'warn',  'Timeout detected',            1, 50);
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.warning',            'warning',            'warn',  'Warning detected',            1, 51);
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.deprecated',         'deprecated',         'warn',  'Deprecation warning',         1, 52);
INSERT OR IGNORE INTO marker_patterns(id, pattern, severity, label, enabled, priority) VALUES
    ('builtin.retry',              'retry',              'warn',  'Retry detected',              1, 53);
