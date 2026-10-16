# Requip User Management System

A comprehensive User Management system built with a Go backend (Gin + GORM + MySQL) and a React frontend (Vite + TypeScript + Tailwind CSS).

## Features
- **Full CRUD**: Create, Read, Update, and Soft Delete users.
- **Pagination**: Paginated user listing on both backend and frontend.
- **Validation**: Strict server-side and client-side validation for fields like Aadhaar, PAN, and Mobile.
- **Modern UI**: Sleek, responsive design using Tailwind CSS v4 and Lucide icons.
- **Unit Testing**: Comprehensive unit tests for backend Service and Handler layers using mocks.

## Tech Stack
- **Backend**: Go 1.25, Gin, GORM, MySQL 8.4.
- **Frontend**: React 19, TypeScript, Vite, Tailwind CSS v4, Axios.
- **Database**: MySQL (running in Docker).

> **Advanced DevOps project.** This repository was taken from a working app (baseline commit [`59da7ed`](https://github.com/BHAV0207/Requip-Assignment-/commit/59da7edbdad65f642915d796a8d33fe806830318), tag `baseline`) to a system that is built, scanned, delivered by GitOps, released by canary, and observable on Kubernetes. See [docs/BASELINE.md](docs/BASELINE.md) for the starting point.

| Document | Contents |
|---|---|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Components, request path, delivery path, failure boundaries |
| [docs/DECISIONS.md](docs/DECISIONS.md) | Each design choice, the alternatives, and when they would be better |
| [docs/RUNBOOK.md](docs/RUNBOOK.md) | What to do for each alert |
| [docs/EVIDENCE.md](docs/EVIDENCE.md) | Pipeline runs, deployed version, user journey, canary and rollback demos |

## Repository layout

```
Backend/                 Go API (+ Dockerfile)
Frontend/                React SPA (+ Dockerfile, nginx.conf)
docker-compose.yml       full stack without Kubernetes
.github/workflows/       CI/CD pipeline
deploy/
  helm/requip/           the application chart (values.yaml + values-kind.yaml)
  argocd/apps/           app-of-apps: every Argo CD Application
  argocd/bootstrap/      Argo CD install values + root Application
  platform/              values for Traefik, Prometheus/Grafana, Loki, Tempo, Alloy, Argo Rollouts, metrics-server
  kind/cluster.yaml      local cluster definition
scripts/                 bootstrap.sh, smoke-test.sh, teardown.sh
docs/                    architecture, decisions, runbook, evidence
```

## Run it on Kubernetes (kind)

### Prerequisites
- Docker (≈ 8 GB RAM and 4 CPUs available to it)
- `kind` ≥ 0.24, `kubectl`, `helm` ≥ 3.14, `openssl`, `git`
- Optional: the `kubectl argo rollouts` plugin, to watch canaries
- Port 80 free on 127.0.0.1

macOS: `brew install kind kubectl helm argoproj/tap/kubectl-argo-rollouts`

### Steps
```bash
git clone https://github.com/BHAV0207/Requip-Assignment-.git && cd Requip-Assignment-
git checkout devops-project      # or main after the PR is merged

./scripts/bootstrap.sh           # ~5-10 min on first run (image pulls)
./scripts/smoke-test.sh          # end-to-end user journey through the ingress
```

`bootstrap.sh` does six things:
1. Creates the kind cluster.
2. Creates the namespaces, with the restricted Pod Security Standard on `requip`.
3. Generates random secrets.
4. Installs Argo CD with Helm.
5. Applies the root Application.
6. Waits until Argo CD has synced the platform and the app.

Everything after Argo CD is installed comes from Git.

If the GHCR images are not public yet, or to test local changes, run `LOCAL_IMAGES=1 ./scripts/bootstrap.sh`. It builds the `:dev` images and loads them into kind.

### Access (bound to 127.0.0.1 only)

| What | URL | Login |
|---|---|---|
| App | http://requip.localtest.me | none |
| Argo CD | http://argocd.localtest.me | `admin` / printed by bootstrap |
| Grafana | http://grafana.localtest.me | `admin` / printed by bootstrap. Dashboard: **Requip — Service Overview** |
| Prometheus | http://prometheus.localtest.me | none |
| Alertmanager | http://alertmanager.localtest.me | none |

`*.localtest.me` is public DNS that always resolves to 127.0.0.1.

Watch a release: `kubectl argo rollouts get rollout requip-backend -n requip --watch`

### Tear down
```bash
./scripts/teardown.sh
```

## How a change reaches users
1. Push to `main` or `devops-project`. A pull request runs the same checks but publishes nothing.
2. GitHub Actions runs:
   - Go tests with the race detector, `govulncheck`
   - ESLint, TypeScript build, `npm audit`
   - `helm lint`, `kubeconform`
   - Trivy scans of the repo (secrets, misconfiguration, CVEs)
3. Each image is built, scanned by Trivy (blocking on fixable HIGH/CRITICAL), given an SBOM, pushed to GHCR as `:<sha7>` and signed with cosign.
4. The `deploy` job (GitHub Environment `kind-local`; add required reviewers there for a manual approval) commits the new tag to `deploy/helm/requip/values-kind.yaml`.
5. Argo CD notices the commit within ~60 s and syncs. Argo Rollouts moves the API to 33%, then 66%, then 100% of pods, while an analysis checks the canary's error rate and p95 latency in Prometheus. A failed check aborts the release and stable pods keep serving.
6. The running version is visible at `/api/v1/version`, in the UI footer, and in the `requip_build_info` metric.

**Rollback:** an automatic abort covers a bad canary. To roll back a release that was promoted, `git revert` its `deploy:` commit and Argo CD deploys the previous tag.

## Run without Kubernetes (Docker Compose)
```bash
cp .env.example .env    # set the two passwords
docker compose up --build
# UI http://localhost:3000   API http://localhost:8080/api/v1/users
```

## Run tests
```bash
cd Backend && go test -race ./...
cd Frontend && npm ci && npx eslint . && npm run build
```

## Known limitations
- **MySQL is a single replica with no backups.** Deleting the PVC or the cluster loses data. Production would use a managed database or operator with backups and replicas.
- **Secrets are created by the bootstrap script, not stored declaratively.** Recreating the cluster generates new passwords. Production: Sealed Secrets or External Secrets Operator.
- **Plain HTTP.** Everything binds to 127.0.0.1 only. Production: TLS at the ingress with cert-manager.
- **Approximate canary weights.** They depend on replica count (no traffic router), see D3.
- **Image signatures are not verified in the cluster.** They are signed in CI, but nothing checks them at admission yet.
- **Alerts go nowhere.** Alertmanager has a null receiver; alerts are visible in the UI only.
- **The API returns full Aadhaar and PAN values** (unchanged app behaviour). They are kept out of logs, but should be masked in responses.
- **Some platform components run without the restricted security context.** Their upstream charts (Loki, Tempo, Prometheus) only apply restricted settings in the `requip` namespace.

## API Documentation

All API endpoints are prefixed with `/api/v1`. By default, the server runs on `http://localhost:8080`.

### User Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/users` | Create a new user record. |
| `GET` | `/users` | List users with pagination using `page` and `limit` query parameters. |
| `GET` | `/users/:id` | Retrieve detailed information for a specific user. |
| `PUT` | `/users/:id` | Update user details (Note: PAN and Aadhaar are immutable). |
| `DELETE` | `/users/:id` | Soft delete a user record. |

#### Create User Request Body
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "primaryMobile": "9876543210",
  "secondaryMobile": "9876543211",
  "aadhaar": "123456789012",
  "pan": "ABCDE1234F",
  "dateOfBirth": "1990-01-01",
  "placeOfBirth": "Mumbai",
  "currentAddress": "123 Street, Mumbai",
  "permanentAddress": "123 Street, Mumbai"
}
```

## Pain Points & Learnings

### Pain Points
- **MySQL Compatibility**: Migrating from Postgres required careful adjustment of UUID generation patterns, as MySQL handles UUIDs differently (using `char(36)` or similar instead of a native `uuid` type with defaults).
- **Port Conflicts**: Encountered issues with the default MySQL port (3306) being used by host services, requiring a shift to 33066 for the containerized database.
- **Tailwind CSS v4 Migration**: Transitioning to Tailwind v4 required a shift from the traditional `tailwind.config.js` to a CSS-first configuration, which has different build-time requirements and linting rules.
- **React 19 Linting**: Stricter `react-hooks/set-state-in-effect` rules in React 19 required refactoring components to use key-based remounting and `useCallback` to avoid cascading renders.

### Learnings
- **Interface-Driven Design**: Extracting interfaces for Repository and Service significantly improved the testability of the backend, allowing for clean mocking without complex setup.
- **Key-Based Remounting**: Using the `key` prop in React to reset complex forms proved to be a much more readable and robust pattern than syncing state via `useEffect`.
- **Phased Development**: Building the project in small, incremental chunks with frequent verification ensured that bugs were caught early and the system remained stable throughput the process.

