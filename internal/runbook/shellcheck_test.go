package runbook

import (
	"testing"
)

func TestValidateShellSyntaxValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command string
	}{
		{"simple echo", "echo hello"},
		{"pipe", "echo hello | grep hello"},
		{"variable", "FOO=bar; echo $FOO"},
		{"if statement", "if true; then echo ok; fi"},
		{"for loop", "for i in 1 2 3; do echo $i; done"},
		{"command substitution", "echo $(date)"},
		{"redirect", "echo hello > /dev/null"},
		{"empty", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			warnings := ValidateShellSyntax(0, tc.command)
			if len(warnings) > 0 {
				t.Errorf("expected no warnings for %q, got: %v", tc.command, warnings)
			}
		})
	}
}

func TestValidateShellSyntaxInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command string
	}{
		{"unclosed quote", "echo 'hello"},
		{"unclosed parenthesis", "echo $(date"},
		{"unclosed brace", "echo ${FOO"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			warnings := ValidateShellSyntax(0, tc.command)
			if len(warnings) == 0 {
				t.Errorf("expected warnings for %q, got none", tc.command)
			}
		})
	}
}

func TestValidateScriptSyntaxValid(t *testing.T) {
	t.Parallel()

	script := `#!/bin/sh
set -e
echo "Starting deploy"
if [ -f /tmp/lock ]; then
    echo "locked"
    exit 1
fi
echo "Done"
`
	warnings := ValidateScriptSyntax(0, script)
	if len(warnings) > 0 {
		t.Errorf("expected no warnings for valid script, got: %v", warnings)
	}
}

func TestValidateScriptSyntaxInvalid(t *testing.T) {
	t.Parallel()

	script := `#!/bin/sh
echo "Starting
if [ -f /tmp/lock ]; then
    echo "locked"
fi
`
	warnings := ValidateScriptSyntax(0, script)
	if len(warnings) == 0 {
		t.Error("expected warnings for script with unclosed quote")
	}
}

func TestValidateRunbookShellSyntax(t *testing.T) {
	t.Parallel()

	steps := []Step{
		{Type: "run", Title: "Valid", Command: "echo ok"},
		{Type: "script", Title: "Invalid script", Script: "echo 'unclosed"},
		{Type: "approval", Title: "Approve", Description: "review"},
		{Type: "run", Title: "Another valid", Command: "ls -la"},
	}

	warnings := ValidateRunbookShellSyntax(steps)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %v", len(warnings), warnings)
	}
	if warnings[0].Step != 1 {
		t.Errorf("warning step = %d, want 1", warnings[0].Step)
	}
}

func TestValidateShellSyntaxFromStrings(t *testing.T) {
	t.Parallel()

	inputs := []ShellCheckInput{
		{Step: 0, Type: "run", Source: "echo ok"},
		{Step: 1, Type: "script", Source: "echo 'unclosed"},
		{Step: 2, Type: "run", Source: "ls -la"},
	}

	warnings := ValidateShellSyntaxFromStrings(inputs)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %v", len(warnings), warnings)
	}
	if warnings[0].Step != 1 {
		t.Errorf("warning step = %d, want 1", warnings[0].Step)
	}
}

func TestShellWarningFields(t *testing.T) {
	t.Parallel()

	warnings := ValidateShellSyntax(3, "echo 'unclosed")
	if len(warnings) == 0 {
		t.Fatal("expected warnings for unclosed quote")
	}

	w := warnings[0]
	if w.Step != 3 {
		t.Errorf("Step = %d, want 3", w.Step)
	}
	if w.Message == "" {
		t.Error("Message should not be empty")
	}
	if w.Line < 1 {
		t.Errorf("Line = %d, want >= 1", w.Line)
	}
	if w.Column < 1 {
		t.Errorf("Column = %d, want >= 1", w.Column)
	}
}

func TestFormatWarnings(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		result := FormatWarnings(nil)
		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("single warning", func(t *testing.T) {
		t.Parallel()
		warnings := []ShellWarning{
			{Step: 0, Line: 1, Column: 5, Message: "unclosed quote"},
		}
		result := FormatWarnings(warnings)
		if result == "" {
			t.Error("expected non-empty string")
		}
	})

	t.Run("multiple warnings", func(t *testing.T) {
		t.Parallel()
		warnings := []ShellWarning{
			{Step: 0, Line: 1, Column: 5, Message: "problem 1"},
			{Step: 2, Line: 3, Column: 1, Message: "problem 2"},
		}
		result := FormatWarnings(warnings)
		if result == "" {
			t.Error("expected non-empty string")
		}
	})
}
