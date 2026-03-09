package store

import (
	"strings"
	"time"
)

func parseStoreTime(raw string) time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}
	}
	if ts, err := time.Parse(time.RFC3339, raw); err == nil {
		return ts.UTC()
	}
	if ts, err := time.Parse("2006-01-02 15:04:05", raw); err == nil {
		return ts.UTC()
	}
	return time.Time{}
}

func parseStoreTimePtr(raw string) *time.Time {
	ts := parseStoreTime(raw)
	if ts.IsZero() {
		return nil
	}
	return &ts
}

func formatTimePtr(ts *time.Time) string {
	if ts == nil || ts.IsZero() {
		return ""
	}
	return ts.UTC().Format(time.RFC3339)
}
