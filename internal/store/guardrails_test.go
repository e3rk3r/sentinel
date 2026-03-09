package store

import (
	"context"
	"testing"
	"time"
)

func TestGuardrailSchemaSeedsDefaultRules(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	defer func() { _ = s.Close() }()

	rules, err := s.ListGuardrailRules(context.Background())
	if err != nil {
		t.Fatalf("ListGuardrailRules: %v", err)
	}
	if len(rules) < 2 {
		t.Fatalf("rules len = %d, want >= 2 defaults", len(rules))
	}

	expected := map[string]struct {
		mode  string
		scope string
	}{
		"action.session.kill.confirm": {mode: GuardrailModeConfirm, scope: GuardrailScopeAction},
		"action.pane.kill.warn":       {mode: GuardrailModeWarn, scope: GuardrailScopeAction},
	}

	found := make(map[string]bool, len(expected))
	for _, rule := range rules {
		exp, ok := expected[rule.ID]
		if !ok {
			continue
		}
		found[rule.ID] = true
		if rule.Mode != exp.mode {
			t.Fatalf("rule %s mode = %q, want %q", rule.ID, rule.Mode, exp.mode)
		}
		if rule.Scope != exp.scope {
			t.Fatalf("rule %s scope = %q, want %q", rule.ID, rule.Scope, exp.scope)
		}
	}
	for id := range expected {
		if !found[id] {
			t.Fatalf("missing seed rule %s", id)
		}
	}
}

func TestGuardrailRuleUpsert(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	if err := s.UpsertGuardrailRule(ctx, GuardrailRuleWrite{
		ID:       "action.window.kill.warn",
		Name:     "Warn on window kill",
		Scope:    GuardrailScopeAction,
		Pattern:  "^window\\.kill$",
		Mode:     GuardrailModeWarn,
		Severity: "info",
		Message:  "Window kill warning",
		Enabled:  true,
		Priority: 30,
	}); err != nil {
		t.Fatalf("UpsertGuardrailRule(create): %v", err)
	}

	if err := s.UpsertGuardrailRule(ctx, GuardrailRuleWrite{
		ID:       "action.window.kill.warn",
		Name:     "Confirm window kill",
		Scope:    GuardrailScopeAction,
		Pattern:  "^window\\.kill$",
		Mode:     GuardrailModeConfirm,
		Severity: "warn",
		Message:  "Window kill confirm",
		Enabled:  true,
		Priority: 12,
	}); err != nil {
		t.Fatalf("UpsertGuardrailRule(update): %v", err)
	}

	rules, err := s.ListGuardrailRules(ctx)
	if err != nil {
		t.Fatalf("ListGuardrailRules: %v", err)
	}
	var found *GuardrailRule
	for i := range rules {
		if rules[i].ID == "action.window.kill.warn" {
			found = &rules[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("expected rule action.window.kill.warn in %+v", rules)
		return
	}
	if found.Mode != GuardrailModeConfirm {
		t.Fatalf("rule mode = %q, want confirm", found.Mode)
	}
	if found.Priority != 12 {
		t.Fatalf("rule priority = %d, want 12", found.Priority)
	}
}

func TestGuardrailScopeNormalizesToAction(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	if err := s.UpsertGuardrailRule(ctx, GuardrailRuleWrite{
		ID:       "test.scope.normalize",
		Name:     "Scope normalization test",
		Scope:    "anything",
		Pattern:  `test`,
		Mode:     GuardrailModeBlock,
		Severity: "error",
		Message:  "blocked",
		Enabled:  true,
		Priority: 50,
	}); err != nil {
		t.Fatalf("UpsertGuardrailRule: %v", err)
	}

	rules, err := s.ListGuardrailRules(ctx)
	if err != nil {
		t.Fatalf("ListGuardrailRules: %v", err)
	}
	for _, r := range rules {
		if r.ID == "test.scope.normalize" {
			if r.Scope != GuardrailScopeAction {
				t.Fatalf("scope = %q, want %q", r.Scope, GuardrailScopeAction)
			}
			return
		}
	}
	t.Fatal("rule not found")
}

func TestGuardrailAuditInsertAndList(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	defer func() { _ = s.Close() }()
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	if _, err := s.InsertGuardrailAudit(ctx, GuardrailAuditWrite{
		RuleID:      "action.session.kill.confirm",
		Decision:    GuardrailModeConfirm,
		Action:      "session.kill",
		SessionName: "dev",
		WindowIndex: -1,
		PaneID:      "",
		Override:    false,
		Reason:      "confirm required",
		MetadataRaw: `{"source":"api"}`,
		CreatedAt:   now,
	}); err != nil {
		t.Fatalf("InsertGuardrailAudit: %v", err)
	}

	rows, err := s.ListGuardrailAudit(ctx, 10)
	if err != nil {
		t.Fatalf("ListGuardrailAudit: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows len = %d, want 1", len(rows))
	}
	if rows[0].Action != "session.kill" || rows[0].Decision != GuardrailModeConfirm {
		t.Fatalf("unexpected audit row: %+v", rows[0])
	}
}
