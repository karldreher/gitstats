package poller

import (
	"strings"
	"testing"
	"time"

	"github.com/karldreher/gitstats/internal/ghclient"
)

func TestParseConventional(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		wantType string
		wantConv string
	}{
		{"feat", "feat: add new feature", "feat", "true"},
		{"fix", "fix: handle nil pointer", "fix", "true"},
		{"fix with scope", "fix(auth): handle nil token", "fix", "true"},
		{"chore", "chore: update deps", "chore", "true"},
		{"docs", "docs: update readme", "docs", "true"},
		{"style", "style: format imports", "style", "true"},
		{"refactor", "refactor: extract helper", "refactor", "true"},
		{"perf", "perf: cache response", "perf", "true"},
		{"test", "test: add unit tests", "test", "true"},
		{"non-conventional", "WIP stuff", "undefined", "false"},
		{"uppercase type", "Feat: something", "undefined", "false"},
		{"missing space after colon", "feat:no space", "undefined", "false"},
		{"unknown type", "build: update makefile", "undefined", "false"},
		{"multiline conventional", "feat: add thing\n\nBody text here", "feat", "true"},
		{"multiline non-conventional", "WIP\n\nsome details", "undefined", "false"},
		// The regex uses .{1,50} without $ — enforces minimum 1 char, not a strict max
		{"description non-empty", "feat: " + strings.Repeat("x", 50), "feat", "true"},
		{"description over 50 still matches", "feat: " + strings.Repeat("x", 51), "feat", "true"},
		{"empty message", "", "undefined", "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotConv := parseConventional(tt.message)
			if gotType != tt.wantType || gotConv != tt.wantConv {
				t.Errorf("parseConventional(%q) = (%q, %q), want (%q, %q)",
					tt.message, gotType, gotConv, tt.wantType, tt.wantConv)
			}
		})
	}
}

func TestPollInterval_Default(t *testing.T) {
	t.Setenv("POLL_INTERVAL_MINUTES", "")
	if got := pollInterval(); got != 15*time.Minute {
		t.Errorf("default interval = %v, want 15m", got)
	}
}

func TestPollInterval_Custom(t *testing.T) {
	t.Setenv("POLL_INTERVAL_MINUTES", "30")
	if got := pollInterval(); got != 30*time.Minute {
		t.Errorf("custom interval = %v, want 30m", got)
	}
}

func TestPollInterval_Invalid(t *testing.T) {
	t.Setenv("POLL_INTERVAL_MINUTES", "notanumber")
	if got := pollInterval(); got != 15*time.Minute {
		t.Errorf("invalid value should fall back to default 15m, got %v", got)
	}
}

func TestResolveSince_PersonalNoStore(t *testing.T) {
	since := resolveSince(ghclient.PersonalMode, nil)
	expected := time.Now().AddDate(0, 0, -30)
	diff := since.Sub(expected)
	if diff < -5*time.Second || diff > 5*time.Second {
		t.Errorf("personal mode without store: want ~30d ago, got diff %v", diff)
	}
}

func TestResolveSince_OrgNoStore(t *testing.T) {
	since := resolveSince(ghclient.OrgMode, nil)
	diff := time.Since(since)
	if diff > 5*time.Second || diff < 0 {
		t.Errorf("org mode without store: want ~now, got %v ago", diff)
	}
}
