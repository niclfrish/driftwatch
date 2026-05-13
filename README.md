# driftwatch

Lightweight daemon that detects config drift between running containers and their source manifests.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git && cd driftwatch && make build
```

---

## Usage

Point driftwatch at your manifests directory and let it run in the background:

```bash
driftwatch --manifests ./k8s --interval 30s
```

It will periodically compare running container configurations against the source manifests and report any drift it finds.

**Example output:**

```
[DRIFT] container=api-server field=image expected=api:v1.2.0 actual=api:v1.1.9
[DRIFT] container=worker field=env.LOG_LEVEL expected=info actual=debug
[OK]    container=db-primary no drift detected
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--manifests` | `./manifests` | Path to source manifest files |
| `--interval` | `60s` | How often to check for drift |
| `--output` | `stdout` | Output format: `stdout`, `json`, or `prometheus` |
| `--kubeconfig` | `~/.kube/config` | Path to kubeconfig file |

---

## Requirements

- Go 1.21+
- Kubernetes cluster or Docker daemon access

---

## License

MIT © 2024 driftwatch contributors