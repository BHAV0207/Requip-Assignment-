# Viva preparation

The viva has three parts: *why this choice*, *how does it behave*, *what would change*. Practise answering each of these out loud with the file open.

## A. Trace a commit to the running pod (have these files open)

1. `git push` → `.github/workflows/ci-cd.yml` runs. Name the jobs: `backend`, `frontend` and `config` run in parallel, then `build` (a matrix of 2), then `deploy`.
2. The image tag is `${GITHUB_SHA::7}`, and the build-args put it in the binary (`Backend/Dockerfile`, `-ldflags -X …version.Version`).
3. Trivy scans the image **before** it is pushed (the `load: true` build, then the scan, then `docker push`). Why that order: a vulnerable image never reaches the registry.
4. The `deploy` job runs `yq` to set the tag in `deploy/helm/requip/values-kind.yaml` and commits with `[skip ci]`. Why `[skip ci]`: otherwise the bot commit would start the pipeline again. (Pushes made with `GITHUB_TOKEN` also don't trigger workflows.)
5. Argo CD polls every 60 s (`timeout.reconciliation` in `deploy/argocd/bootstrap/argocd-values.yaml`). It renders the chart with `values.yaml` + `values-kind.yaml`, sees that the Rollout pod template changed, and applies it.
6. Argo Rollouts creates a new ReplicaSet, and the steps in `values.yaml` take it 33% → pause 60 s → 66% → pause 60 s → 100%. From step 1, a background AnalysisRun (`templates/backend-analysis.yaml`) queries Prometheus every 30 s.
7. Verify with `curl http://requip.localtest.me/api/v1/version`, the UI footer, and `kubectl argo rollouts get rollout requip-backend -n requip`.

## B. Trace a user request

Browser → `requip.localtest.me` (DNS resolves to 127.0.0.1) → kind port mapping 80 → NodePort 30080 (`deploy/kind/cluster.yaml`) → Traefik → Ingress `templates/ingress.yaml` (`/api` goes to `requip-backend:8080`, `/` to `requip-frontend:80`) → Service → one ready pod → Gin middleware chain in `cmd/server/main.go` (otel → metrics → logging → recovery → CORS) → handler → service → repository → `otelsql` → MySQL headless Service → `requip-mysql-0`.

NetworkPolicies allowed each hop (`templates/networkpolicy.yaml`). Name which policy allows each one.

## C. Likely "why" questions, with the short answer

| Question | Answer (expand from DECISIONS.md) |
|---|---|
| Why a StatefulSet for MySQL? | Stable identity plus its own PVC. A Deployment could start two pods on one RWO volume during an update (D2) |
| Why not a StatefulSet for the API? | It's stateless; pods are interchangeable. A StatefulSet would only slow rollouts down (ordered, one pod at a time) |
| Why canary rather than blue-green? | Exposes a small share of users while measuring real traffic, and needs no double capacity. Blue-green suits versions that cannot run side by side (D3) |
| Why does liveness not check the DB? | A DB outage would restart every API pod for no gain. Readiness handles it instead (D9) |
| Why no CPU limit? | CFS throttling hurts latency; requests already guarantee each pod its share (D10) |
| Why GitOps instead of `kubectl apply` from CI? | Cloud runners can't reach a laptop cluster. Also: audit trail, drift correction, no cluster credentials in CI (D4) |
| Why are secrets not in Git? | Base64 is not encryption. Created out-of-band; Sealed Secrets or ESO in production (D11) |
| Why Traefik? | ingress-nginx was retired upstream in March 2026 (D16) |
| Why multi-window burn-rate alerts? | Fast detection of big outages, no paging for blips, quick reset after recovery (D15) |
| Why omit `replicas` from the Rollout? | The HPA owns it; otherwise Argo CD and the HPA fight (D8) |
| Why replace the gorm OTel plugin? | It pulled in Postgres and ClickHouse drivers with critical CVEs that the app never uses (D13) |

## D. "What would change if…" (practise these out loud)

- **Traffic grows 100×.** Change these, in this order:
  - The HPA scales the API, but MySQL is the bottleneck: add read replicas, a connection pooler, caching.
  - Raise `DB_MAX_OPEN_CONNS` carefully. (replicas × pool) must stay below MySQL `max_connections`.
  - Watch the *DB connection pool* panel (`waits/s`).
- **The DB must survive a node loss.** local-path PVs are node-bound. Use a networked StorageClass (EBS/PD) or a managed DB, plus backups (e.g., a `mysqldump` CronJob to object storage), and *test the restore*.
- **A new release adds a NOT NULL column.** Startup AutoMigrate plus a canary means the old pods break. Use expand/contract migrations in a PreSync Job: add the column as nullable → deploy → backfill → add the constraint (D17).
- **Prometheus is down during a rollout.** The analysis measurement errors and the rollout aborts. It fails *safe*.
- **The image tag in Git doesn't exist in the registry.** New pods hit ImagePullBackOff, the canary never becomes ready, and stable keeps serving. Fix it in Git.
- **Someone runs `kubectl scale deploy requip-frontend --replicas=0`.** Argo CD self-heal puts it back to 2 (measured: spec restored in ~6 s, pods ready in ~11 s, see evidence/09). Demo this live: `kubectl -n requip scale deploy requip-frontend --replicas=0 && kubectl -n requip get deploy requip-frontend -w`.
- **You need a 1% canary or header-based routing.** Add a traffic router (Traefik weighted services, a service mesh, or Gateway API) to `strategy.canary.trafficRouting` (D3).
- **Release needs human approval.** Add required reviewers to the `kind-local` GitHub Environment, or add an indefinite `pause: {}` step and promote with `kubectl argo rollouts promote`.

## E. Live demos you can run in the viva (each takes about 2–5 minutes)

```bash
# 1. user journey
./scripts/smoke-test.sh

# 2. self-healing pod
kubectl -n requip delete pod -l app.kubernetes.io/component=api --wait=false; kubectl -n requip get pods -w

# 3. database failure (see evidence/06): API goes NotReady, 0 restarts, data kept
kubectl -n requip delete pod requip-mysql-0

# 4. drift correction
kubectl -n requip scale deploy requip-frontend --replicas=0

# 5. bad release (edit values-kind.yaml: backend.env.FAULT_ERROR_RATE: "0.5", push)
kubectl argo rollouts get rollout requip-backend -n requip --watch
```

## F. Weak spots to admit before they're found

- A single MySQL with no backups is the largest risk to the system.
- Users saw about 10 s of 500s during the DB test. Readiness can't help when every pod shares the dependency.
- The canary weight is approximate (replica-based).
- Images are signed but nothing verifies the signatures at admission.
- The create path makes 4 uniqueness `SELECT`s, then an INSERT (seen in the trace). That's 5 round trips and a check-then-insert race; the unique indexes are the real guarantee.
- The API returns full Aadhaar and PAN values; only logs and traces are protected.
