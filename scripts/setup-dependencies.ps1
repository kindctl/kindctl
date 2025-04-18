# Check if Docker is installed and running
Write-Host "Checking for Docker..."
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "Error: Docker is not installed. Please install Docker Desktop or Rancher Desktop."
    exit 1
}
try {
    docker ps | Out-Null
    Write-Host "Docker is installed and running."
}
catch {
    Write-Host "Error: Docker is not running. Please start Docker Desktop or Rancher Desktop."
    exit 1
}

# Function to add to PATH
function Add-ToPath($path) {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", [System.EnvironmentVariableTarget]::User)
    if ($currentPath -notlike "*$path*") {
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$path", [System.EnvironmentVariableTarget]::User)
        $env:Path += ";$path"
    }
}

# Install winget if not available (for Windows 10/11)
if (-not (Get-Command winget -ErrorAction SilentlyContinue)) {
    Write-Host "Installing winget..."
    Invoke-WebRequest -Uri https://github.com/microsoft/winget-cli/releases/latest/download/Microsoft.DesktopAppInstaller_8wekyb3d8bbwe.msixbundle -OutFile winget.msixbundle
    Add-AppxPackage -Path winget.msixbundle
    Remove-Item winget.msixbundle
}

# Install kind
if (-not (Get-Command kind -ErrorAction SilentlyContinue)) {
    Write-Host "Installing kind..."
    Invoke-WebRequest -Uri https://kind.sigs.k8s.io/dl/v0.23.0/kind-windows-amd64 -OutFile kind.exe
    New-Item -ItemType Directory -Path "$env:ProgramFiles\kind" -Force
    Move-Item -Path kind.exe -Destination "$env:ProgramFiles\kind\kind.exe"
    Add-ToPath "$env:ProgramFiles\kind"
}
else {
    Write-Host "kind is already installed."
}

# Install kubectl
if (-not (Get-Command kubectl -ErrorAction SilentlyContinue)) {
    Write-Host "Installing kubectl..."
    $kubectlVersion = (Invoke-WebRequest -Uri https://dl.k8s.io/release/stable.txt -UseBasicParsing).Content.Trim()
    Invoke-WebRequest -Uri "https://dl.k8s.io/release/$kubectlVersion/bin/windows/amd64/kubectl.exe" -OutFile kubectl.exe
    New-Item -ItemType Directory -Path "$env:ProgramFiles\kubectl" -Force
    Move-Item -Path kubectl.exe -Destination "$env:ProgramFiles\kubectl\kubectl.exe"
    Add-ToPath "$env:ProgramFiles\kubectl"
}
else {
    Write-Host "kubectl is already installed."
}

# Install helm
if (-not (Get-Command helm -ErrorAction SilentlyContinue)) {
    Write-Host "Installing helm..."
    Invoke-WebRequest -Uri https://get.helm.sh/helm-v3.15.4-windows-amd64.zip -OutFile helm.zip
    Expand-Archive -Path helm.zip -DestinationPath helm-temp
    New-Item -ItemType Directory -Path "$env:ProgramFiles\helm" -Force
    Move-Item -Path helm-temp\windows-amd64\helm.exe -Destination "$env:ProgramFiles\helm\helm.exe"
    Add-ToPath "$env:ProgramFiles\helm"
    Remove-Item -Recurse helm.zip, helm-temp
}
else {
    Write-Host "helm is already installed."
}

# Verify installations
Write-Host "Verifying installations..."
foreach ($cmd in @("kind", "kubectl", "helm")) {
    if (Get-Command $cmd -ErrorAction SilentlyContinue) {
        Write-Host "$cmd installed: $((Invoke-Expression "$cmd version --client") -join ' ')"
    }
    else {
        Write-Host "Error: $cmd installation failed."
        exit 1
    }
}

Write-Host "All dependencies installed successfully! You can now use kindctl."