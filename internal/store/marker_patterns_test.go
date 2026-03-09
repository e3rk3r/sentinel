package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func TestMarkerPatternSeedsDefaults(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	defer func() { _ = s.Close() }()

	patterns, err := s.ListMarkerPatterns(context.Background())
	if err != nil {
		t.Fatalf("ListMarkerPatterns: %v", err)
	}
	if len(patterns) < 8 {
		t.Fatalf("patterns len = %d, want >= 8 defaults", len(patterns))
	}

	expectedIDs := map[string]string{
		"builtin.panic":              "error",
		"builtin.fatal":              "error",
		"builtin.error":              "error",
		"builtin.oom":                "error",
		"builtin.segfault":           "error",
		"builtin.segmentation-fault": "error",
		"builtin.connection-refused": "error",
		"builtin.killed":             "error",
		"builtin.timeout":            "warn",
		"builtin.warning":            "warn",
		"builtin.deprecated":         "warn",
		"builtin.retry":              "warn",
	}

	found := make(map[string]bool, len(expectedIDs))
	for _, p := range patterns {
		wantSeverity, ok := expectedIDs[p.ID]
		if !ok {
			continue
		}
		found[p.ID] = true
		if p.Severity != wantSeverity {
			t.Fatalf("pattern %s severity = %q, want %q", p.ID, p.Severity, wantSeverity)
		}
		if !p.Enabled {
			t.Fatalf("pattern %s should be enabled by default", p.ID)
		}
	}
	for id := range expectedIDs {
		if !found[id] {
			t.Fatalf("missing seed pattern %s", id)
		}
	}
}

func TestMarkerPatternSeedsOrderedByPriority(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	defer func() { _ = s.Close() }()

	patterns, err := s.ListMarkerPatterns(context.Background())
	if err != nil {
		t.Fatalf("ListMarkerPatterns: %v", err)
	}
	for i := 1; i < len(patterns); i++ {
		if patterns[i].Priority < patterns[i-1].Priority {
			t.Fatalf("patterns not ordered by priority: %d (%s) < %d (%s)",
				patterns[i].Priority, patterns[i].ID,
				patterns[i-1].Priority, patterns[i-1].ID)
		}
	}
}

func TestMarkerPatternUpsertAndList(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	// Insert a new pattern.
	if err := s.UpsertMarkerPattern(ctx, MarkerPatternWrite{
		ID:       "custom.deploy-failed",
		Pattern:  "deploy failed",
		Severity: "error",
		Label:    "Deployment failure",
		Enabled:  true,
		Priority: 5,
	}); err != nil {
		t.Fatalf("UpsertMarkerPattern(create): %v", err)
	}

	// Update it.
	if err := s.UpsertMarkerPattern(ctx, MarkerPatternWrite{
		ID:       "custom.deploy-failed",
		Pattern:  "deploy error",
		Severity: "warn",
		Label:    "Deploy error",
		Enabled:  false,
		Priority: 99,
	}); err != nil {
		t.Fatalf("UpsertMarkerPattern(update): %v", err)
	}

	patterns, err := s.ListMarkerPatterns(ctx)
	if err != nil {
		t.Fatalf("ListMarkerPatterns: %v", err)
	}

	var found *MarkerPattern
	for i := range patterns {
		if patterns[i].ID == "custom.deploy-failed" {
			found = &patterns[i]
			break
		}
	}
	if found == nil {
		t.Fatal("custom.deploy-failed not found after upsert")
	}
	if found.Pattern != "deploy error" {
		t.Fatalf("pattern = %q, want %q", found.Pattern, "deploy error")
	}
	if found.Severity != "warn" {
		t.Fatalf("severity = %q, want %q", found.Severity, "warn")
	}
	if found.Label != "Deploy error" {
		t.Fatalf("label = %q, want %q", found.Label, "Deploy error")
	}
	if found.Enabled {
		t.Fatal("expected disabled after update")
	}
	if found.Priority != 99 {
		t.Fatalf("priority = %d, want 99", found.Priority)
	}
}

func TestMarkerPatternDelete(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	// Delete an existing seed pattern.
	if err := s.DeleteMarkerPattern(ctx, "builtin.panic"); err != nil {
		t.Fatalf("DeleteMarkerPattern: %v", err)
	}

	patterns, err := s.ListMarkerPatterns(ctx)
	if err != nil {
		t.Fatalf("ListMarkerPatterns: %v", err)
	}
	for _, p := range patterns {
		if p.ID == "builtin.panic" {
			t.Fatal("builtin.panic should have been deleted")
		}
	}
}

func TestMarkerPatternDeleteNonexistent(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	err := s.DeleteMarkerPattern(ctx, "nonexistent")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("DeleteMarkerPattern(nonexistent) err = %v, want sql.ErrNoRows", err)
	}
}

func TestMarkerPatternUpsertValidation(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	t.Run("empty id", func(t *testing.T) {
		t.Parallel()
		err := s.UpsertMarkerPattern(ctx, MarkerPatternWrite{
			ID:      "",
			Pattern: "test",
		})
		if err == nil {
			t.Fatal("expected error for empty id")
		}
	})

	t.Run("empty pattern", func(t *testing.T) {
		t.Parallel()
		err := s.UpsertMarkerPattern(ctx, MarkerPatternWrite{
			ID:      "test",
			Pattern: "",
		})
		if err == nil {
			t.Fatal("expected error for empty pattern")
		}
	})
}

func TestMarkerPatternSeverityNormalization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"error", "error"},
		{"ERROR", "error"},
		{"info", "info"},
		{"INFO", "info"},
		{"warn", "warn"},
		{"WARN", "warn"},
		{"anything", "warn"},
		{"", "warn"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			s := newTestStore(t)
			defer func() { _ = s.Close() }()
			ctx := context.Background()

			id := "test.severity." + tt.input
			if tt.input == "" {
				id = "test.severity.empty"
			}
			if err := s.UpsertMarkerPattern(ctx, MarkerPatternWrite{
				ID:       id,
				Pattern:  "test",
				Severity: tt.input,
				Enabled:  true,
			}); err != nil {
				t.Fatalf("UpsertMarkerPattern: %v", err)
			}

			patterns, err := s.ListMarkerPatterns(ctx)
			if err != nil {
				t.Fatalf("ListMarkerPatterns: %v", err)
			}
			for _, p := range patterns {
				if p.ID == id {
					if p.Severity != tt.want {
						t.Fatalf("severity = %q, want %q", p.Severity, tt.want)
					}
					return
				}
			}
			t.Fatalf("pattern %s not found", id)
		})
	}
}
