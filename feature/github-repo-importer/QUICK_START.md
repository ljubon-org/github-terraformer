# Quick Start - GitHub Repo Importer

## Setup & Usage

See the main documentation: [LOCAL_DEVELOPMENT_SETUP.md](../../docs/LOCAL_DEVELOPMENT_SETUP.md)

## Quick Commands

```bash
# Setup (one-time)
source .env

# Import single repo
just import-repo owner/repo

# Import all repos
just import-repos

# Compare configs
just compare dir1 dir2

# Run tests
just test
```

## Common Patterns

```bash
# Import with environments enabled
echo "feature_github_environment: true" >> import-config.yaml
just import-repo owner/repo

# Custom output directory
go run main.go import owner/repo -o ../custom-dir

# Debug mode
export DEBUG=1
go run main.go import owner/repo
```

For full setup instructions and troubleshooting, see [LOCAL_DEVELOPMENT_SETUP.md](../../docs/LOCAL_DEVELOPMENT_SETUP.md)