# kindctl

`kindctl` is a CLI tool to set up and manage local Kubernetes clusters using Kind for development purposes. It allows you to configure tools like Postgres, Redis, and Kubernetes Dashboard via a YAML file and automatically manages ingress and `/etc/hosts` entries.

## Prerequisites

- **Docker Desktop** or **Rancher Desktop**: Ensure Docker is installed and running.
    - Verify: `docker ps`
- The `setup-dependencies` script will install the following if not present:
    - `kind`: For creating Kubernetes clusters.
    - `kubectl`: For interacting with the cluster.
    - `helm`: For installing chart-based tools (e.g., Postgres, Redis).

## Installation

### Setup Dependencies

Before installing `kindctl`, run the dependency setup script to ensure `kind`, `kubectl`, and `helm` are installed.

#### Linux/macOS

```bash
curl -L https://github.com/<your-username>/kindctl/raw/main/scripts/setup-dependencies.sh -o setup-dependencies.sh
chmod +x setup-dependencies.sh
./setup-dependencies.sh
```

#### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://github.com/<your-username>/kindctl/raw/main/scripts/setup-dependencies.ps1 -OutFile setup-dependencies.ps1
.\setup-dependencies.ps1
```

### Install kindctl

#### Linux/macOS

```bash
curl -L https://github.com/<your-username>/kindctl/releases/latest/download/install.sh | bash
```

#### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://github.com/<your-username>/kindctl/releases/latest/download/install.ps1 -OutFile install.ps1
.\install.ps1
```

## Usage

1. **Initialize a project**:

   ```bash
   kindctl init
   ```

   This creates a `kindctl.yaml` file with a default Kubernetes Dashboard configuration, sets up a Kind cluster, and installs the NGINX ingress controller.

2. **Update the cluster**:

   ```bash
   kindctl update
   ```

   Reads `kindctl.yaml`, installs enabled tools, and automatically sets up required Helm repositories for tools like Postgres, Redis, pgAdmin, and RabbitMQ.

3. **Destroy the cluster**:

   ```bash
   kindctl destroy
   ```

   Deletes the Kind cluster specified in `kindctl.yaml`.

## Configuration

Example `kindctl.yaml`:

```yaml
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
```

## Supported Tools

- Kubernetes Dashboard
- Postgres
- Redis
- pgAdmin
- Adminer
- RabbitMQ
- Mailpit

See the GitHub Pages for detailed tool configurations.

## Building from Source

```bash
git clone https://github.com/<your-username>/kindctl
cd kindctl
./scripts/build.sh
```

## License

MIT