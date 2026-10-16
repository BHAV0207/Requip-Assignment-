# Decision record

Each entry: **what was chosen, why it fits this app, what else was considered, and when the other option would be better.**

---

### D1. Local kind cluster (1 control plane + 2 workers)
- **Why:** the brief allows local Kubernetes, and it costs nothing. Two workers make pod spreading, PodDisruptionBudgets and "a node died" scenarios real. A single-node Minikube cannot show them.
- **Alternatives:** Minikube (one node by default), k3d, a managed cloud cluster (EKS/GKE).
- **Choose differently when:** real users or uptime matter. Then a managed cluster with a cloud load balancer, managed storage classes and a managed database.

### D2. Rollout (canary) for the API; Deployment for the SPA; StatefulSet for MySQL
- **API:** stateless and the part most likely to break in a release, so it gets progressive delivery.
- **SPA:** static files with nothing to measure per pod, so a rolling update with `maxUnavailable: 0` is enough.
- **MySQL:** needs a stable name (`requip-mysql-0`) and its own PVC that follows it, which is exactly what a StatefulSet provides. A Deployment with a PVC would allow two pods to fight over one RWO volume during a rollout.
- **Choose differently when:** in production, use a managed database (RDS/Cloud SQL) or an operator (e.g., MySQL Operator) for backups, replication and failover. A single-replica StatefulSet has none of these.

### D3. Canary without a traffic router (replica-weighted)
- **How it works:** Argo Rollouts scales the new ReplicaSet to 1 of 3 pods (≈33%) and then 2 of 3 (≈66%). The single Service balances across all pods, so traffic follows the pod ratio. An AnalysisRun queries Prometheus every 30 s for **the canary's own pods** (filtered by `rollouts_pod_template_hash`), checking error ratio ≤ 5% and p95 ≤ 500 ms. One failure aborts the rollout.
- **Why:** it works with any ingress, adds no extra moving parts, and is easy to explain.
- **Trade-off:** weights are coarse (they depend on the replica count) and a user can hit different versions on consecutive requests.
- **Alternatives:**
  - *Rolling update* (no automated safety net)
  - *Blue-green* (needs double capacity; instant switch and instant rollback; good when the two versions can't run side by side)
  - *Traffic-router canary* (Traefik/Istio/Gateway API: exact percentages, header-based routing)
  - *Feature flags* (decouple deploy from release for individual features)
- **Choose differently when:** you need precise small percentages (1%) or sticky sessions. Then use a traffic router.

### D4. GitOps with Argo CD (pull) instead of `kubectl apply`/`helm upgrade` from CI (push)
- **Why:** GitHub's cloud runners cannot reach a cluster on a laptop. With pull-based GitOps the cluster fetches its desired state instead. It also gives:
  - an audit trail (every deploy is a Git commit)
  - drift correction (`selfHeal` reverts `kubectl edit`)
  - rollback by `git revert`
  - no cluster credentials stored in CI
- **Alternatives:** a self-hosted runner with a kubeconfig plus `helm upgrade` (simpler, but CI then holds cluster-admin credentials and drift goes unnoticed); Flux (equivalent to Argo CD, without the UI).
- **Choose differently when:** a team has no Kubernetes expertise or the target isn't Kubernetes. Then a classic push pipeline is simpler.

### D5. CI writes the image tag into Git (instead of Argo CD Image Updater)
- **Why:** it is explicit and easy to trace. The `deploy: requip <sha7>` commit links to the CI run that built the image. The `deploy` job runs in the `kind-local` GitHub Environment, so adding *required reviewers* there turns it into a manual approval gate without changing the workflow.
- **Trade-off:** bot commits appear in the history.
- **Alternative:** Argo CD Image Updater watches the registry, which means fewer commits but less explicit traceability.

### D6. Immutable image tags = short commit SHA
- **Why:** `latest` is not traceable and doesn't trigger a rollout. A SHA tag maps one-to-one to source code. The binary also embeds the full SHA (`-ldflags`), exposed at `/api/v1/version` and in the `requip_build_info` metric.
- **Further hardening:** deploy by digest and verify cosign signatures at admission (Kyverno/Sigstore policy-controller). The images are already signed, but nothing in the cluster checks the signatures yet.

### D7. Helm, one chart with per-environment values files
- **Why:** one source of truth for the manifests. Environment differences (`values-kind.yaml`) are small overrides, not copied YAML. Helm also handles conditionals: Rollout vs Deployment, HPA on or off.
- **Alternative:** Kustomize overlays (no templating language; good when the differences are patches rather than logic).

### D8. HPA on CPU, and the Rollout omits `replicas` when HPA is on
- **Why:** if Git specified `replicas: 3` and the HPA scaled to 5, Argo CD would see drift and scale back to 3. The two controllers would fight. Leaving the field out (plus `ignoreDifferences` on `/spec/replicas`) makes the HPA the only owner.
- **Choose differently when:** load is driven by queue depth or requests per second, not CPU. Then KEDA or custom metrics.

### D9. Probes
- **Liveness `/healthz`:** checks only that the process is alive. It deliberately does **not** check the DB. A DB outage would otherwise restart every API pod, turning one failure into two.
- **Readiness `/readyz`:** pings MySQL with a 2 s timeout, so pods that can't serve are taken out of rotation.
- **Startup probe:** allows 90 s, covering the app's DB-connect retry loop, before liveness starts counting.
- **`preStop: sleep 5` plus graceful shutdown on SIGTERM:** endpoints are removed before the server stops accepting connections, so rolling updates drop no requests.

### D10. Resource requests and a memory limit, but no CPU limit
- Requests let the scheduler place pods and give the HPA a baseline. A memory limit protects the node from a leak (the pod is OOM-killed and restarted).
- No CPU limit, because CFS throttling inflates tail latency even when the node has idle CPU. Requests already guarantee each pod its share.

### D11. Secrets created out-of-band, not stored in Git
- `scripts/bootstrap.sh` generates random passwords with `openssl rand` and creates the Kubernetes Secrets directly. The chart only refers to them through `existingSecret`. Nothing secret lives in Git, CI or the images.
- **Trade-off:** the secret isn't declarative, so recreating the cluster generates new passwords (fine for a lab whose data is disposable).
- **Choose differently when:** there are multiple environments or the cluster is long-lived. Then use Sealed Secrets (encrypted in Git) or External Secrets Operator with Vault or a cloud secret manager.

### D12. Defence in depth at runtime
- Restricted Pod Security Standard enforced on the namespace: non-root, no privilege escalation, all capabilities dropped, seccomp `RuntimeDefault`, read-only root filesystem.
- Distroless API image: no shell or package manager, 54 MB.
- `automountServiceAccountToken: false`: the app never talks to the Kubernetes API, so it gets no token to steal.
- Default-deny NetworkPolicies (see ARCHITECTURE.md). MySQL is reachable only from API pods.
- Argo CD `AppProject requip`: the app can only deploy from this repo, into the `requip` namespace, with **no** cluster-scoped resources.

### D13. Security checks in CI
| Check | Catches |
|---|---|
| `govulncheck` | Known Go vulnerabilities **in code paths actually called** |
| `npm audit --omit=dev` | Vulnerable runtime JS dependencies |
| `trivy fs` (secret, misconfig, vuln) | Committed secrets, Dockerfile/Kubernetes misconfigurations, lockfile CVEs |
| `trivy image` | OS/library CVEs in the final image; the push happens only after this passes |
| `kubeconform` | Invalid manifests before they reach the cluster |
| SBOM (SPDX) + cosign keyless signing | Supply-chain record and provenance of each image |

Every GitHub Action is pinned to a full commit SHA, and the workflow defaults to `contents: read`.

These checks earned their place while building: they found 8 vulnerable npm packages, a Dockerfile without `USER`, and a tracing plugin that pulled in unused Postgres and ClickHouse drivers with critical CVEs. That plugin was replaced with `otelsql`.

### D14. Observability stack: Prometheus + Grafana + Loki + Tempo + Alloy
- **Why:** open source, runs locally, and Grafana links all three signals (metric → log → trace).
- **Alloy:** reads logs through the Kubernetes API, so it needs no privileged hostPath mount on nodes.
- **Loki and Tempo:** run single-binary with local storage. Production would use object storage.
- **Alternatives:** the ELK/EFK stack (heavier); a hosted service like Grafana Cloud or Datadog (no ops burden, but costs money).

### D15. SLOs and burn-rate alerts
- **SLO:**
  - **99.5%** of `/api` requests succeed (not 5xx)
  - p95 latency under 300 ms
- **Alerts:** multi-window burn-rate (fast 5m and 1h at 14.4×, slow 30m and 6h at 6×) instead of "error rate > X%". The alert fires quickly for a serious outage, stays quiet for a single blip, and resolves soon after recovery.
- **Response:** every alert links to a runbook section (`docs/RUNBOOK.md`).

### D16. Traefik instead of ingress-nginx
- ingress-nginx was retired upstream in March 2026 and no longer receives security fixes. Traefik is maintained, supports the standard `Ingress` API (so the chart stays portable) and also supports the Gateway API.

### D17. Schema migration at API startup (GORM AutoMigrate)
- **Kept from the original app,** and only additive changes are expected.
- **Risk:** several replicas migrating at once, or a canary changing the schema under the stable version.
- **Choose differently when:** migrations become non-additive. Then run a separate migration Job (an Argo CD `PreSync` hook) using expand/contract migrations, and set `RUN_MIGRATIONS=false` on the API.
