# kindctl

`kindctl` is a CLI tool to set up and manage local Kubernetes clusters using Kind for development purposes. It allows you to configure tools like Postgres, Redis, and Kubernetes Dashboard via a YAML file and automatically manages ingress and `/etc/hosts` entries.

## Installation

### Linux/macOS

```bash
curl -L https://github.com/<your-username>/kindctl/releases/latest/download/install.sh | bash
```

### Windows

```powershell
Invoke-WebRequest -Uri https://github.com/<your-username>/kindctl/releases/latest/download/install.ps1 -OutFile install.ps1; .\install.ps1
```

## Usage

1. **Initialize a project**:

```bash
   kindctl init
   ```

This creates a `kindctl.yaml` file with a default Kubernetes Dashboard configuration and sets up a Kind cluster.

2. **Update the cluster**: Edit `kindctl.yaml` to enable tools, then run:

```bash
   kindctl update
   ```

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

MIT# kindctl
