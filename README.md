# gitstats

A GitHub commit metrics exporter. Polls GitHub for commit activity and exposes Prometheus metrics labeled by repo, author, commit type, and conventional commit compliance.

## Modes

**Personal** — tracks one user's commits via `contributionsCollection`. Set `GITHUB_PAT` and `GITHUB_USER`. Performs a 30-day backfill on first run.

**Org** — tracks all non-archived repos in a GitHub org via GitHub App auth. Set `GITHUB_APP_ID`, `GITHUB_APP_INSTALLATION_ID`, `GITHUB_APP_PRIVATE_KEY`, and `GITHUB_ORG`.

## Environment Variables

| Variable | Description |
|---|---|
| `POLL_INTERVAL_MINUTES` | Polling interval (default: `15`) |
| `GITHUB_PAT` | Personal access token (personal mode) |
| `GITHUB_USER` | GitHub login to track (personal mode) |
| `GITHUB_APP_ID` | GitHub App ID (org mode) |
| `GITHUB_APP_INSTALLATION_ID` | Installation ID (org mode) |
| `GITHUB_APP_PRIVATE_KEY` | PEM-encoded RSA private key (org mode) |
| `GITHUB_ORG` | Org login name (org mode) |
| `PERSISTENCE_FILE` | Path to JSON state file (optional) |
| `PERSISTENCE_REDIS_HOST` | Redis host:port (optional) |
| `PERSISTENCE_REDIS_PASS` | Redis password (required if host is set) |

## Endpoints

| Endpoint | Description |
|---|---|
| `GET /metrics` | Prometheus scrape target |
| `GET /healthz` | Liveness probe (always 200) |
| `GET /readyz` | Readiness probe (200 after first poll completes) |

## Metrics

`gitstats_commits_total` — counter with labels `repo`, `author`, `commit_type`, `conventional`.

```promql
# Conventional commit compliance ratio
sum(rate(gitstats_commits_total{conventional="true"}[1h])) /
sum(rate(gitstats_commits_total[1h]))

# Per-author compliance
sum by (author) (rate(gitstats_commits_total{conventional="true"}[1h]))
```

## Running

```bash
docker compose up --build
```
