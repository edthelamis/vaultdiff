# vaultdiff

> A CLI tool to diff and audit changes between Vault secret versions across environments.

---

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultdiff.git
cd vaultdiff
go build -o vaultdiff .
```

---

## Usage

Compare two versions of a secret within the same Vault path:

```bash
vaultdiff --path secret/myapp/config --v1 3 --v2 5
```

Diff secrets across environments:

```bash
vaultdiff --path secret/myapp/config \
  --env-a staging \
  --env-b production
```

Audit all changes to a secret over time:

```bash
vaultdiff audit --path secret/myapp/config --since 2024-01-01
```

### Flags

| Flag | Description |
|------|-------------|
| `--path` | Vault secret path |
| `--v1` | First version to compare |
| `--v2` | Second version to compare |
| `--env-a` | Source environment |
| `--env-b` | Target environment |
| `--output` | Output format: `text`, `json`, `yaml` (default: `text`) |

---

## Requirements

- Go 1.21+
- HashiCorp Vault with KV v2 secrets engine
- `VAULT_ADDR` and `VAULT_TOKEN` environment variables set

---

## License

MIT © 2024 Your Name