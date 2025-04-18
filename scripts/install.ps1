$Version = "v0.1.0"
$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") { "amd64" } else { "arm64" }
$Binary = "kindctl-windows-$Arch.exe"
$Url = "https://github.com/kindctl/kindctl/releases/download/$Version/$Binary"
$InstallPath = "$env:ProgramFiles\kindctl"

Write-Host "Installing kindctl $Version for Windows/$Arch..."

Invoke-WebRequest -Uri $Url -OutFile kindctl.exe
New-Item -ItemType Directory -Path $InstallPath -Force
Move-Item -Path kindctl.exe -Destination "$InstallPath\kindctl.exe"
$env:Path += ";$InstallPath"
[Environment]::SetEnvironmentVariable("Path", $env:Path, [System.EnvironmentVariableTarget]::User)

Write-Host "kindctl installed successfully!"