# Personal Mode

Personal mode authenticates as a GitHub user via a Personal Access Token (PAT) and tracks that user's contributions using the GraphQL `contributionsCollection` API.

## Creating a Personal Access Token

Use a fine-grained PAT. You can use either pre-filled link below to open GitHub's token creation page with the name and description already filled in, then set an expiration and click **Generate token**.

**Public repos only:**
```
https://github.com/settings/personal-access-tokens/new?name=gitstats&description=Track+GitHub+contribution+statistics
```

**Including private repos** (after clicking, scope the token to the relevant repositories):
```
https://github.com/settings/personal-access-tokens/new?name=gitstats&description=Track+GitHub+contribution+statistics&contents=read
```

Or create one manually:

1. Go to **github.com → Settings → Developer settings → Personal access tokens → Fine-grained tokens**
2. Click **Generate new token**
3. Set a descriptive name (e.g. `gitstats`) and an expiration appropriate for your use case
4. Under **Permissions**, set only what is required:

| Permission | Level | Required | Reason |
|---|---|---|---|
| (none) | — | Public repos only | `contributionsCollection` for public repos is accessible to any authenticated token |
| Contents | Read-only | Private repos | Required to include contributions from private repositories |

5. Click **Generate token** and copy the value immediately

> **Note:** Fine-grained PATs are best suited for local development and small self-hosted deployments. Regardless of use case, rotate tokens regularly and use the minimum required scopes.

## Configuration

Set in `.env`:

```
GITHUB_PAT=ghp_...
GITHUB_USER=your-github-login
```

---

# Org Mode

Org mode uses a GitHub App to authenticate and tracks commits across all non-archived repositories in an organization.

## Creating a GitHub App

1. Go to **github.com → Your Organization → Settings → Developer settings → GitHub Apps → New GitHub App**
2. Fill in the required fields:
   - **GitHub App name**: any name (e.g. `gitstats`)
   - **Homepage URL**: any URL (required by GitHub; not used by this app)
3. Under **Webhook**, uncheck **Active** — webhooks are not used
4. Under **Repository permissions**, set only:

| Permission | Level | Reason |
|---|---|---|
| Metadata | Read-only | Required for all GitHub Apps; needed to list repositories |
| Contents | Read-only | Required to read commits from repositories |

5. Leave all other permissions at **No access**
6. Under **Where can this GitHub App be installed?**, select **Only on this account**
7. Click **Create GitHub App**
8. Note the **App ID** displayed at the top of the app settings page

## Generating a Private Key

1. On the app settings page, scroll to **Private keys**
2. Click **Generate a private key**
3. A `.pem` file is downloaded — store it securely

## Installing the App

1. On the app settings page, click **Install App** in the left sidebar
2. Click **Install** next to your organization
3. Choose **All repositories** (or select specific ones)
4. Click **Install**
5. After installation, the URL will be of the form:
   ```
   https://github.com/organizations/<org>/settings/installations/<installation-id>
   ```
   Note the **Installation ID** from the URL

## Configuration

Set in `.env`:

```
GITHUB_APP_ID=123456
GITHUB_APP_INSTALLATION_ID=78901234
GITHUB_APP_PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----
MIIEow...
-----END RSA PRIVATE KEY-----"
GITHUB_ORG=your-org-name
```

> **Note:** `GITHUB_APP_PRIVATE_KEY` must contain the full PEM content with literal newlines preserved. In a `.env` file, wrap the value in double quotes.
