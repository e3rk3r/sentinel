package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

// MarkerPattern represents a configurable marker detection pattern used by
// the watchtower pane content scanner to detect timeline markers.
type MarkerPattern struct {
	ID        string    `json:"id"`
	Pattern   string    `json:"pattern"`
	Severity  string    `json:"severity"`
	Label     string    `json:"label"`
	Enabled   bool      `json:"enabled"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"createdAt"`
}

// MarkerPatternWrite is the input type for UpsertMarkerPattern.
type MarkerPatternWrite struct {
	ID       string
	Pattern  string
	Severity string
	Label    string
	Enabled  bool
	Priority int
}

// ListMarkerPatterns returns all marker patterns ordered by priority then id.
func (s *Store) ListMarkerPatterns(ctx context.Context) ([]MarkerPattern, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, pattern, severity, label, enabled, priority, created_at
		   FROM marker_patterns
		  ORDER BY priority ASC, id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make([]MarkerPattern, 0, 16)
	for rows.Next() {
		var (
			row          MarkerPattern
			enabledRaw   int
			createdAtRaw string
		)
		if err := rows.Scan(
			&row.ID,
			&row.Pattern,
			&row.Severity,
			&row.Label,
			&enabledRaw,
			&row.Priority,
			&createdAtRaw,
		); err != nil {
			return nil, err
		}
		row.Enabled = enabledRaw == 1
		row.CreatedAt = parseStoreTime(createdAtRaw)
		out = append(out, row)
	}
	return out, rows.Err()
}

// UpsertMarkerPattern creates or updates a marker pattern.
func (s *Store) UpsertMarkerPattern(ctx context.Context, row MarkerPatternWrite) error {
	id := strings.TrimSpace(row.ID)
	if id == "" {
		return errors.New("marker pattern id is required")
	}
	pattern := strings.TrimSpace(row.Pattern)
	if pattern == "" {
		return errors.New("pattern is required")
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO marker_patterns(id, pattern, severity, label, enabled, priority, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, datetime('now'))
		 ON CONFLICT(id) DO UPDATE SET
			pattern  = excluded.pattern,
			severity = excluded.severity,
			label    = excluded.label,
			enabled  = excluded.enabled,
			priority = excluded.priority`,
		id,
		pattern,
		normalizeMarkerSeverity(row.Severity),
		strings.TrimSpace(row.Label),
		boolToInt(row.Enabled),
		row.Priority,
	)
	return err
}

// DeleteMarkerPattern removes a marker pattern by id.
func (s *Store) DeleteMarkerPattern(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("marker pattern id is required")
	}
	result, err := s.db.ExecContext(ctx, `DELETE FROM marker_patterns WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func normalizeMarkerSeverity(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "error":
		return "error"
	case timelineSeverityInfo:
		return timelineSeverityInfo
	default:
		return "warn"
	}
}
