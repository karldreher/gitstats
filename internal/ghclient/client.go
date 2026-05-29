// Package ghclient provides a GitHub API client supporting both personal (PAT) and org (GitHub App) auth modes.
package ghclient

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Mode distinguishes personal (PAT) from org (GitHub App) operation.
type Mode int

const (
	PersonalMode Mode = iota
	OrgMode
)

// CommitData is a normalized commit record returned by either fetch mode.
type CommitData struct {
	Message string // full commit message
	Repo    string // "owner/repo"
	Author  string // GitHub login or display name
}

// Client handles GitHub API authentication and dispatch.
type Client struct {
	httpClient     *http.Client
	mode           Mode
	user           string
	token          string
	org            string
	appID          string
	installationID string
	privateKey     *rsa.PrivateKey
	tokenExpiresAt time.Time
}

// FromEnv constructs a Client from environment variables.
// Exactly one of GITHUB_PAT (personal) or GITHUB_APP_ID (org) must be set.
func FromEnv() (*Client, error) {
	pat := os.Getenv("GITHUB_PAT")
	appID := os.Getenv("GITHUB_APP_ID")

	if pat != "" && appID != "" {
		return nil, errors.New("only one GitHub auth mode may be configured: unset GITHUB_PAT or GITHUB_APP_ID")
	}
	if pat == "" && appID == "" {
		return nil, errors.New("GitHub auth required: set GITHUB_PAT (personal mode) or GITHUB_APP_ID (org mode)")
	}

	c := &Client{httpClient: &http.Client{Timeout: 30 * time.Second}}

	if pat != "" {
		user := os.Getenv("GITHUB_USER")
		if user == "" {
			return nil, errors.New("GITHUB_USER is required when GITHUB_PAT is set")
		}
		c.mode = PersonalMode
		c.token = pat
		c.user = user
		return c, nil
	}

	installationID := os.Getenv("GITHUB_APP_INSTALLATION_ID")
	org := os.Getenv("GITHUB_ORG")
	cert := os.Getenv("GITHUB_APP_PRIVATE_KEY")

	if installationID == "" || org == "" || cert == "" {
		return nil, errors.New("GITHUB_APP_INSTALLATION_ID, GITHUB_ORG, and GITHUB_APP_PRIVATE_KEY are all required in org mode")
	}

	cert = strings.ReplaceAll(cert, `\n`, "\n")
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cert))
	if err != nil {
		return nil, fmt.Errorf("parsing GITHUB_APP_PRIVATE_KEY: %w", err)
	}

	c.mode = OrgMode
	c.org = org
	c.appID = appID
	c.installationID = installationID
	c.privateKey = key

	if err := c.refreshInstallationToken(); err != nil {
		return nil, fmt.Errorf("getting installation token: %w", err)
	}
	return c, nil
}

// Mode returns the operation mode (PersonalMode or OrgMode).
func (c *Client) Mode() Mode { return c.mode }

// User returns the GitHub login or org name the client is authenticated as.
func (c *Client) User() string { return c.user }

// FetchCommits dispatches to the appropriate mode's fetch implementation.
func (c *Client) FetchCommits(since time.Time) ([]CommitData, error) {
	switch c.mode {
	case PersonalMode:
		return c.fetchPersonalCommits(since)
	case OrgMode:
		return c.fetchOrgCommits(since)
	}
	return nil, errors.New("unknown mode")
}

// bearerToken returns the current auth token, refreshing if needed for org mode.
func (c *Client) bearerToken() (string, error) {
	if c.mode == PersonalMode {
		return c.token, nil
	}
	if time.Now().Add(5 * time.Minute).After(c.tokenExpiresAt) {
		if err := c.refreshInstallationToken(); err != nil {
			return "", err
		}
	}
	return c.token, nil
}

func (c *Client) refreshInstallationToken() error {
	now := time.Now()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": now.Add(-60 * time.Second).Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"iss": c.appID,
	})
	signed, err := tok.SignedString(c.privateKey)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", c.installationID)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("installation token request failed: %s", resp.Status)
	}

	var result struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	expiresAt, err := time.Parse(time.RFC3339, result.ExpiresAt)
	if err != nil {
		return err
	}
	c.token = result.Token
	c.tokenExpiresAt = expiresAt
	return nil
}
