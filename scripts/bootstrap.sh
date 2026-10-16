#!/usr/bin/env bash
# Creates the local kind cluster and hands everything else to Argo CD.
#
#   ./scripts/bootstrap.sh                 # track the current Git branch
#   REVISION=main ./scripts/bootstrap.sh   # track a specific branch
#   LOCAL_IMAGES=1 ./scripts/bootstrap.sh  # also build :dev images and load them into kind
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CLUSTER=requip
REVISION="${REVISION:-$(git -C "$ROOT" rev-parse --abbrev-ref HEAD)}"
ARGOCD_CHART_VERSION=10.10.5

need() { command -v "$1" >/dev/null || { echo "missing required tool: $1" >&2; exit 1; }; }
for t in docker kind kubectl helm openssl; do need "$t"; done

echo "==> 1/6 kind cluster"
if ! kind get clusters | grep -qx "$CLUSTER"; then
  kind create cluster --config "$ROOT/deploy/kind/cluster.yaml" --wait 120s
fi
kubectl config use-context "kind-$CLUSTER" >/dev/null

echo "==> 2/6 namespaces (Pod Security 'restricted' enforced on the app namespace)"
kubectl create namespace requip --dry-run=client -o yaml | kubectl apply -f -
kubectl label namespace requip --overwrite \
  pod-security.kubernetes.io/enforce=restricted \
  pod-security.kubernetes.io/warn=restricted
kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -

echo "==> 3/6 secrets (generated locally, never stored in Git)"
if ! kubectl -n requip get secret requip-db >/dev/null 2>&1; then
  kubectl -n requip create secret generic requip-db \
    --from-literal=mysql-root-password="$(openssl rand -hex 16)" \
    --from-literal=mysql-password="$(openssl rand -hex 16)"
fi
if ! kubectl -n monitoring get secret grafana-admin >/dev/null 2>&1; then
  kubectl -n monitoring create secret generic grafana-admin \
    --from-literal=admin-user=admin \
    --from-literal=admin-password="$(openssl rand -hex 12)"
fi

if [[ "${LOCAL_IMAGES:-0}" == "1" ]]; then
  echo "==> optional: build :dev images and load them into kind"
  docker build -t ghcr.io/bhav0207/requip-backend:dev --build-arg VERSION=dev "$ROOT/Backend"
  docker build -t ghcr.io/bhav0207/requip-frontend:dev "$ROOT/Frontend"
  kind load docker-image --name "$CLUSTER" ghcr.io/bhav0207/requip-backend:dev ghcr.io/bhav0207/requip-frontend:dev
fi

echo "==> 4/6 Argo CD $ARGOCD_CHART_VERSION"
helm repo add argo https://argoproj.github.io/argo-helm >/dev/null 2>&1 || true
helm upgrade --install argocd argo/argo-cd --version "$ARGOCD_CHART_VERSION" \
  --namespace argocd --create-namespace \
  -f "$ROOT/deploy/argocd/bootstrap/argocd-values.yaml" --wait --timeout 10m

echo "==> 5/6 root application (tracking '$REVISION')"
sed "s/__REVISION__/$REVISION/g" "$ROOT/deploy/argocd/bootstrap/root-app.yaml" | kubectl apply -f -

echo "==> 6/6 waiting for Argo CD to sync everything (first run pulls many images)"
for i in $(seq 1 90); do
  status=$(kubectl -n argocd get application requip -o jsonpath='{.status.sync.status}/{.status.health.status}' 2>/dev/null || true)
  [[ "$status" == "Synced/Healthy" ]] && break
  printf '\r    requip app: %-25s (%ss)' "${status:-pending}" "$((i * 10))"
  sleep 10
done
echo
kubectl -n argocd get applications

cat <<MSG

Done. Open (bound to 127.0.0.1 only):
  App         http://requip.localtest.me
  Argo CD     http://argocd.localtest.me      admin / $(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d)
  Grafana     http://grafana.localtest.me     admin / $(kubectl -n monitoring get secret grafana-admin -o jsonpath='{.data.admin-password}' | base64 -d)
  Prometheus  http://prometheus.localtest.me
MSG
