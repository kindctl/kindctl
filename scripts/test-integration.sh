#!/bin/bash

set -e

# Ensure dependencies are installed
echo "Setting up dependencies..."
chmod +x scripts/setup-dependencies.sh
./scripts/setup-dependencies.sh

# Build kindctl
echo "Building kindctl..."
chmod +x scripts/build.sh
./scripts/build.sh
KINDCTL_BIN="bin/kindctl-$(uname -s | tr '[:upper:]' '[:lower:]')-amd64"
chmod +x "$KINDCTL_BIN"

# Create test directory
TEST_DIR="test-kindctl"
rm -rf "$TEST_DIR"
mkdir "$TEST_DIR"
cd "$TEST_DIR"

# Test kindctl init
echo "Testing kindctl init..."
"$OLDPWD/$KINDCTL_BIN" init
if [ ! -f "kindctl.yaml" ]; then
    echo "Error: kindctl.yaml not created"
    exit 1
fi
if ! kind get clusters | grep -q "kind-cluster"; then
    echo "Error: Kind cluster not created"
    exit 1
fi
echo "kindctl init passed"

# Create test kindctl.yaml with multiple tools
cat <<EOF > kindctl.yaml
logging:
  level: debug
cluster:
  name: kind-cluster
dashboard:
  enabled: true
  ingress: dashboard.local
postgres:
  enabled: true
  ingress: postgres.local
  version: "16"
  username: postgres
  password: postgres
  database: postgres
redis:
  enabled: true
  ingress: redis.local
EOF

# Test kindctl update
echo "Testing kindctl update..."
"$OLDPWD/$KINDCTL_BIN" update
if ! kubectl get pods -n default | grep -q "postgres"; then
    echo "Error: Postgres pod not found"
    exit 1
fi
if ! kubectl get pods -n default | grep -q "redis"; then
    echo "Error: Redis pod not found"
    exit 1
fi
if ! kubectl get ingress -n default | grep -q "dashboard-ingress"; then
    echo "Error: Dashboard ingress not found"
    exit 1
fi
echo "kindctl update passed"

# Test kindctl destroy
echo "Testing kindctl destroy..."
"$OLDPWD/$KINDCTL_BIN" destroy
if kind get clusters | grep -q "kind-cluster"; then
    echo "Error: Kind cluster not deleted"
    exit 1
fi
echo "kindctl destroy passed"

echo "All integration tests passed!"
cd ..
rm -rf "$TEST_DIR"