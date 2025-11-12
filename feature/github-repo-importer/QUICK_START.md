# Quick Start - GitHub Repo Importer

## Setup (One-Time)

```bash
# 1. Navigate to importer directory
cd github-terraformer/feature/github-repo-importer

# 2. Create .env file
cat > .env << 'EOF'
export GITHUB_TOKEN="ghp_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export OWNER="your-org-name"
EOF

# 3. Link required config files
ln -sf ../../../gcss-config-repo/config/import-config.yaml import-config.yaml
ln -sf ../../../gcss-config-repo/config/app-list.yaml app-list.yaml

# Or create new one:
cat > import-config.yaml << 'EOF'
ignored_repos:
  - "your-org/gcss-config-repo"
  - "your-org/github-terraformer"
feature_github_environment: true
EOF

# 4. Test setup
source .env
go run main.go --help
```

## Daily Usage

```bash
# Always start with
source .env

# Import single repository
just import-repo your-org/repo-name
# OR
go run main.go import your-org/repo-name

# Import all repositories
just import-repos
# OR
go run main.go bulk-import

# Compare configurations
just compare dir1 dir2
# OR
go run main.go compare dir1 dir2
```

## Common Tasks

### Import with Environments

```bash
# Enable in config
echo "feature_github_environment: true" >> import-config.yaml

# Import
go run main.go import your-org/repo-name
```

### Import Specific Repos Only

```bash
cat > import-config.yaml << 'EOF'
selected_repos:
  - "your-org/repo1"
  - "your-org/repo2"
EOF

go run main.go bulk-import
```

### Custom Output Directory

```bash
# Single import
go run main.go import your-org/repo -o ../output-dir

# Bulk import
go run main.go bulk-import -o ../output-dir
```

## Testing & Development

```bash
# Run tests
just test
# OR
go test ./...

# Test specific package
go test ./pkg/github -v

# Build binary
go build -o github-importer

# Format code
go fmt ./...
```

## Output Locations

Default output directories:
- **Single import**: Current directory or specified with `-o`
- **Bulk import**: `importer_tmp_dir/` or specified with `-o`
- **JSON dumps**: `jsondumps/` (for debugging)

## Troubleshooting

```bash
# Token not set
source .env

# Check token validity
curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/user

# Debug mode
export DEBUG=1
go run main.go import your-org/repo

# Check rate limit
curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/rate_limit
```

## File Structure

```
Output YAML structure:
importer_tmp_dir/
├── repo1.yaml        # Repository configuration
├── repo2.yaml        # Repository configuration
└── ...

JSON dumps (debug):
jsondumps/
├── repo.json         # Raw API response
├── environments.json # Environment details
└── rulesets.json     # Ruleset details
```