# Runbook

Each alert in `deploy/helm/requip/templates/monitoring.yaml` links to a section here.
The approach is the same every time: **symptom → hypothesis → evidence → change → verify**.

Useful views:
- Grafana → dashboard **Requip — Service Overview** (http://grafana.localtest.me)
- `kubectl argo rollouts get rollout requip-backend -n requip --watch`
- `kubectl -n argocd get applications`

## high-error-rate

**Alert:** `RequipErrorBudgetFastBurn` (page) or `RequipErrorBudgetSlowBurn` (ticket). More than 14.4× (or 6×) the allowed 0.5% of `/api` requests are returning 5xx.

1. **Did a release just happen?** Look at the *Running versions* panel and `kubectl argo rollouts get rollout requip-backend -n requip`. If a canary is running, analysis normally aborts it on its own. To abort by hand: `kubectl argo rollouts abort requip-backend -n requip`.
2. **Is it one pod or all of them?** *Traffic by pod template* and `sum by (pod)(rate(http_requests_total{status=~"5.."}[5m]))`.
3. **What do the errors say?** In the dashboard's log panel, or in Grafana Explore → Loki: `{namespace="requip", container="api"} | json | level="ERROR"`. Click a `trace_id` to open the trace in Tempo and see which span failed (often the MySQL query).
4. **Is the database healthy?** See [database-down](#database-down).
5. **Roll back a bad version permanently:** `git revert` the `deploy: requip <tag>` commit (or set `backend.image.tag` back to the previous tag in `deploy/helm/requip/values-kind.yaml`) and push. Argo CD syncs the old version. Don't `kubectl edit` in the cluster: self-heal reverts manual changes.
6. **Verify recovery:** `requip:api_error_ratio:rate5m` drops below 0.005, the alert resolves, and `./scripts/smoke-test.sh` passes.

## high-latency

**Alert:** `RequipHighLatency`. API p95 latency has been above 300 ms for 10 minutes.

1. *Latency by route*: one route or all of them?
2. *DB connection pool*: if `waits/s` is above 0, requests are queueing for a connection. Raise `DB_MAX_OPEN_CONNS` in `values.yaml` or find the slow query in Tempo.
3. CPU saturation: `kubectl top pods -n requip`. Check whether the HPA has already scaled out (`kubectl get hpa -n requip`).
4. Verify recovery: p95 back under 300 ms on the dashboard.

## backend-down

**Alert:** `RequipBackendDown`. Prometheus cannot scrape any backend pod.

1. `kubectl get pods -n requip -l app.kubernetes.io/component=api`. Look for CrashLoopBackOff, ImagePullBackOff or Pending.
2. `kubectl describe pod <pod>` and `kubectl logs <pod> --previous`.
   - ImagePullBackOff: the tag in Git does not exist in GHCR, or the package is private.
   - CrashLoop with `database connection failed`: see [database-down](#database-down).
3. If the NetworkPolicy was changed, check that the monitoring namespace can still reach port 8080.

## database-down

**Alert:** `RequipDatabaseNotReady`. The MySQL StatefulSet has no ready replica. API pods go **NotReady** (their `/readyz` pings MySQL) and are taken out of the Service. Liveness keeps passing, so they are **not** restarted.

1. `kubectl get pod requip-mysql-0 -n requip` and `kubectl describe` it.
2. `kubectl get pvc -n requip`: the PVC must be `Bound`.
3. `kubectl logs requip-mysql-0 -n requip`
4. Once MySQL is ready again, API pods become Ready by themselves within ~10 s and need no restart.
5. Verify with `./scripts/smoke-test.sh`.

Data survives pod deletion because it lives on the PVC `data-requip-mysql-0`. It does **not** survive deleting the PVC or the kind cluster: there is no backup in this lab setup (see the README's limitations section).
