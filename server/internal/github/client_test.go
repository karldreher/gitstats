package github

import (
	"testing"
)

func TestFromEnv_NoAuth(t *testing.T) {
	t.Setenv("GITHUB_PAT", "")
	t.Setenv("GITHUB_APP_ID", "")
	_, err := FromEnv()
	if err == nil {
		t.Error("no auth configured: expected error")
	}
}

func TestFromEnv_BothModes(t *testing.T) {
	t.Setenv("GITHUB_PAT", "ghp_test")
	t.Setenv("GITHUB_APP_ID", "12345")
	_, err := FromEnv()
	if err == nil {
		t.Error("both auth modes set: expected error")
	}
}

func TestFromEnv_PersonalMissingUser(t *testing.T) {
	t.Setenv("GITHUB_PAT", "ghp_test")
	t.Setenv("GITHUB_APP_ID", "")
	t.Setenv("GITHUB_USER", "")
	_, err := FromEnv()
	if err == nil {
		t.Error("PAT without GITHUB_USER: expected error")
	}
}

func TestFromEnv_PersonalMode(t *testing.T) {
	t.Setenv("GITHUB_PAT", "ghp_test")
	t.Setenv("GITHUB_APP_ID", "")
	t.Setenv("GITHUB_USER", "testuser")
	client, err := FromEnv()
	if err != nil {
		t.Fatalf("valid personal config: unexpected error: %v", err)
	}
	if client.Mode() != PersonalMode {
		t.Errorf("mode = %v, want PersonalMode", client.Mode())
	}
}

func TestFromEnv_OrgMissingFields(t *testing.T) {
	t.Setenv("GITHUB_PAT", "")
	t.Setenv("GITHUB_APP_ID", "12345")
	t.Setenv("GITHUB_APP_INSTALLATION_ID", "")
	t.Setenv("GITHUB_ORG", "")
	t.Setenv("GITHUB_APP_PRIVATE_KEY", "")
	_, err := FromEnv()
	if err == nil {
		t.Error("org mode with missing fields: expected error")
	}
}
