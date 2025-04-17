# Ensure dependencies are installed
Write-Host "Setting up dependencies..."
.\scripts\setup-dependencies.ps1

# Build kindctl
Write-Host "Building kindctl..."
.\scripts\build.sh
$KINDCTL_BIN = "bin\kindctl-windows-amd64.exe"
if (-not (Test-Path $KINDCTL_BIN)) {
    Write-Host "Error: kindctl binary not found"
    exit 1
}

# Create test directory
$TEST_DIR = "test-kindctl"
Remove-Item -Recurse -Force $TEST_DIR -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Path $TEST_DIR
Set-Location $TEST_DIR

# Test kindctl init
Write-Host "Testing kindctl init..."
& "..\${KINDCTL_BIN}" init
if (-not (Test-Path "kindctl.yaml")) {
    Write-Host "Error: kindctl.yaml not created"
    exit 1
}
$clusters = kind get clusters
if ($clusters -notlike "*kind-cluster*") {
    Write-Host "Error: Kind cluster not created"
    exit 1
}
Write-Host "kindctl init passed"

# Create test kindctl.yaml with multiple tools
$CONFIG_CONTENT = @"
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
"@
Set-Content -Path "kindctl.yaml" -Value $CONFIG_CONTENT

# Test kindctl update
Write-Host "Testing kindctl update..."
& "..\${KINDCTL_BIN}" update
$pods = kubectl get pods -n default
if ($pods -notlike "*postgres*") {
    Write-Host "Error: Postgres pod not found"
    exit 1
}
if ($pods -notlike "*redis*") {
    Write-Host "Error: Redis pod not found"
    exit 1
}
$ingress = kubectl get ingress -n default
if ($ingress -notlike "*dashboard-ingress*") {
    Write-Host "Error: Dashboard ingress not found"
    exit 1
}
Write-Host "kindctl update passed"

# Test kindctl destroy
Write-Host "Testing kindctl destroy..."
& "..\${KINDCTL_BIN}" destroy
$clusters = kind get clusters
if ($clusters -like "*kind-cluster*") {
    Write-Host "Error: Kind cluster not deleted"
    exit 1
}
Write-Host "kindctl destroy passed"

Write-Host "All integration tests passed!"
Set-Location ..
Remove-Item -Recurse -Force $TEST_DIR