# Kubernetes Deployment

Example manifests for deploying gitstats to a Kubernetes cluster.

## Prerequisites

- `kubectl` configured against your target cluster
- A GitHub Personal Access Token **or** a GitHub App (see [Authentication](#authentication))

## Quick Start

1. Create the namespace:
   ```sh
   kubectl create namespace gitstats
   ```

2. Edit `secret.yaml` — fill in your credentials for the auth mode you chose (see [Authentication](#authentication) below).

3. Optionally edit `configmap.yaml` to set `POLL_INTERVAL_MINUTES` or enable a persistence backend.

4. Apply all manifests:
   ```sh
   kubectl apply -f examples/k8s/
   ```

5. Verify the pod is running and ready:
   ```sh
   kubectl -n gitstats get pods
   ```
   The pod starts as `0/1 Ready` until the first poll completes, then transitions to `1/1 Ready`.

6. Confirm metrics are being exported:
   ```sh
   kubectl -n gitstats port-forward svc/gitstats 8000:8000
   curl localhost:8000/readyz    # 200 OK after first poll
   curl localhost:8000/metrics   # gitstats_commits{...} counters
   ```

## Authentication

gitstats uses a GitHub App to track all non-archived repositories in your organization. Edit `secret.yaml` and fill in the four values:

```yaml
GITHUB_APP_ID: "<app-id>"
GITHUB_APP_INSTALLATION_ID: "<installation-id>"
GITHUB_APP_PRIVATE_KEY: |
  -----BEGIN RSA PRIVATE KEY-----
  <paste PEM key here>
  -----END RSA PRIVATE KEY-----
GITHUB_ORG: "<org-name>"
```

The GitHub App needs **Contents: Read** and **Metadata: Read** permissions, installed on the target organization. You can find the App ID on the App's settings page; the Installation ID is in the URL when you view the App's installation under your org's settings (`/organizations/<org>/settings/installations/<id>`).

## Persistence

Without persistence, counter state is lost on pod restart. Choose one backend:

### File (default for simple setups)

1. Uncomment `PERSISTENCE_FILE` in `configmap.yaml`.
2. Apply `pvc.yaml` for storage:
   ```sh
   kubectl apply -f examples/k8s/pvc.yaml
   ```
3. Uncomment the `volumeMounts` and `volumes` blocks in `deployment.yaml`.

### Redis

1. Uncomment `PERSISTENCE_REDIS_HOST` in `configmap.yaml` and set it to your Redis `host:port`.
2. Uncomment `PERSISTENCE_REDIS_PASS` in `secret.yaml` and set the password.
3. No volume or PVC is needed.

## Prometheus Integration

### Annotation-based scraping

`service.yaml` includes `prometheus.io/scrape` annotations. This works with a standard Prometheus deployment configured to scrape annotated services.

### Prometheus Operator (kube-prometheus-stack)

Apply `servicemonitor.yaml` after ensuring the `release` label matches your Prometheus Operator instance:

```sh
kubectl apply -f examples/k8s/servicemonitor.yaml
```

The default label is `release: prometheus`. Verify with:
```sh
kubectl get prometheus -n monitoring -o jsonpath='{.items[*].spec.serviceMonitorSelector}'
```

## Scaling Note

**Do not scale beyond `replicas: 1`.** gitstats uses cumulative Prometheus counters backed by a single persistence store. Multiple replicas would cause each instance to independently increment the same counters, resulting in double-counted metrics.
