# Local Development Setup for GitHub Repo Importer

This guide explains how to set up your local environment for developing and using the GitHub Repo Importer CLI tool.

## Overview

The GitHub Repo Importer is a Go-based CLI tool that imports GitHub repository configurations into YAML files for Terraform management.

## Prerequisites

1. **Go** (1.19+)
   ```bash
   # Check version
   go version

   # Install if needed: https://golang.org/doc/install
   ```

2. **Just** (command runner)
   ```bash
   # Install on Ubuntu/Debian
   curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin

   # Or with cargo
   cargo install just

   # Or download from: https://github.com/casey/just/releases
   ```

3. **GitHub Access**
   - Personal Access Token or GitHub App credentials
   - Required scopes: `repo`, `read:org`, `admin:org`

## Directory Structure

```
github-repo-importer/
├── cmd/                      # Cobra CLI commands
│   ├── root.go              # Root command setup
│   ├── import.go            # Import single repository
│   ├── bulkImport.go        # Import multiple repositories
│   └── compare.go           # Compare configurations
├── pkg/
│   ├── github/              # GitHub API interactions
│   │   ├── github.go        # Main GitHub client
│   │   ├── repositories.go  # Repository structures
│   │   └── constants.go     # Feature flags and constants
│   ├── file/                # File operations
│   │   └── file.go         # YAML file handling
│   └── compare/            # Configuration comparison
├── main.go                  # Entry point
├── Justfile                # Task automation
├── go.mod                  # Go modules
└── import-config.yaml      # Import configuration (local)
```

## Step 1: Environment Setup

Create a `.env` file in the importer directory:

```bash
cd github-terraformer/feature/github-repo-importer

cat > .env << 'EOF'
# GitHub Authentication
export GITHUB_TOKEN="ghp_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export OWNER="your-org-name"

# Feature Flags (via import-config.yaml)
# These are controlled in import-config.yaml, not env vars
EOF

```

## Step 2: Import Configuration

The importer uses `import-config.yaml` for configuration. Create or link it:

### Option A: Link from gcss-config-repo (Recommended)

```bash
# Create symlink to shared config
ln -sf ../../../gcss-config-repo/config/import-config.yaml import-config.yaml

# Create symlink for app list (needed for GitHub App bypass actors)
ln -sf ../../../gcss-config-repo/config/app-list.yaml app-list.yaml

# Optional: Create symlink for output directory (if you want imports to go directly to gcss-config-repo)
ln -sf ../../../gcss-config-repo/importer_tmp_dir importer_tmp_dir
```
EOF
```

## Step 3: Build and Test

```bash
# Load environment
source .env

# Run tests
just test

# Or with go directly
go test ./...

# Build the binary
go build -o github-importer

# Or run directly without building
go run main.go --help
```

## Step 4: Using the Importer

### Import Single Repository

```bash
# Using just
just import-repo your-org/repository-name

# Using go run
go run main.go import your-org/repository-name

# With custom output directory
go run main.go import your-org/repository-name -o ../custom-output-dir

# With custom config
go run main.go import your-org/repository-name -c custom-config.yaml
```

### Bulk Import Repositories

```bash
# Using just
just import-repos

# Using go run
go run main.go bulk-import

# With custom config
go run main.go bulk-import -c import-config.yaml

# With custom output directory
go run main.go bulk-import -o ../output-dir
```

### Compare Configurations

```bash
# Using just
just compare dirA dirB

# Using go run
go run main.go compare dirA dirB

# Compare with verbose output
go run main.go compare dirA dirB -v
```

## Step 5: Output Structure

The importer generates YAML files in the output directory:

```yaml
# output-dir/repository-name.yaml
description: "Repository description"
homepage_url: "https://example.com"
visibility: private
default_branch: main
has_issues: true
has_projects: true
has_wiki: false
has_downloads: true
allow_merge_commit: true
allow_rebase_merge: true
allow_squash_merge: true
delete_branch_on_merge: true
vulnerability_alerts_enabled: true

# If feature_github_environment: true
environments:
  - environment: production
    wait_timer: 300
    can_admins_bypass: false
    prevent_self_review: true
    reviewers:
      users:
        - username1
      teams:
        - team-slug
    # IMPORTANT: deployment_ref_policy uses EITHER protected_branches_policy
    # OR selected_branches_or_tags_policy, but NOT both
    deployment_ref_policy:
      # Option 1: Only protected branches can deploy
      protected_branches_policy: true

  - environment: staging
    deployment_ref_policy:
      # Option 2: Custom branch/tag patterns (requires protected_branches_policy: false or omitted)
      protected_branches_policy: false
      selected_branches_or_tags_policy:
        branch_patterns:
          - "release/*"
          - "hotfix/*"
        tag_patterns:
          - "v*"

# Rulesets (if present)
rulesets:
  - name: "Main Branch Protection"
    target: "branch"
    enforcement: "active"
    # ... ruleset configuration
```

## Development Workflow

### 1. Making Changes

```bash
# Create a feature branch
git checkout -b feature/my-improvement

# Make changes to code
vim pkg/github/github.go

# Run tests
just test

# Test specific functionality
go test ./pkg/github -v

# Test import with your changes
go run main.go import your-org/test-repo
```

### 2. Testing Import Features

```bash
# Test environment import
echo "feature_github_environment: true" >> import-config.yaml
go run main.go import your-org/repo-with-environments

# Test ruleset import
go run main.go import your-org/repo-with-rulesets

# Test bulk import with filters
cat > test-config.yaml << EOF
selected_repos:
  - "your-org/test-repo1"
  - "your-org/test-repo2"
feature_github_environment: true
EOF
go run main.go bulk-import -c test-config.yaml
```

### 3. Debugging

```bash
# Enable debug output
export DEBUG=1
go run main.go import your-org/repo

# Use delve debugger
go get -u github.com/go-delve/delve/cmd/dlv
dlv debug main.go -- import your-org/repo

# Check generated JSON dumps
ls -la jsondumps/
```

## Common Commands Reference

```bash
# Source environment
source .env

# Run tests
just test
go test ./...
go test ./pkg/github -v -run TestSpecificFunction

# Import operations
just import-repo owner/repo
just import-repos
just compare dir1 dir2

# Direct go commands
go run main.go import owner/repo
go run main.go bulk-import
go run main.go compare dir1 dir2

# Build
go build -o github-importer
./github-importer import owner/repo

# Format code
go fmt ./...
gofmt -w .

# Lint
golangci-lint run

# Dependencies
go mod tidy
go mod download
```

## Troubleshooting

### Issue: "GITHUB_TOKEN not set"

**Solution**: Set and export the token:
```bash
export GITHUB_TOKEN="ghp_your_token_here"
# Or source .env file
source .env
```

### Issue: "401 Unauthorized"

**Cause**: Invalid or expired GitHub token

**Solution**: Generate a new token with required permissions:
- Go to GitHub Settings → Developer settings → Personal access tokens
- Required scopes: `repo`, `read:org`, `admin:org`

### Issue: "403 rate limit exceeded"

**Solution**: Wait for rate limit reset or use GitHub App authentication:
```bash
# Check rate limit
curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/rate_limit
```

### Issue: "repository not found"

**Causes**:
1. Repository doesn't exist
2. No access to private repository
3. Wrong organization name

**Solution**: Verify repository exists and token has access:
```bash
gh repo view your-org/repo-name
```

### Issue: Feature flags not working

**Solution**: Ensure import-config.yaml is in the correct location:
```bash
# Default locations checked:
# 1. ./import-config.yaml
# 2. Specified with -c flag
ls -la import-config.yaml
cat import-config.yaml | grep feature_
```

## Adding New Features

1. **Add feature flag** in `pkg/github/constants.go`:
   ```go
   const FeatureMyNewFeature = "feature_my_new_feature"
   ```

2. **Read feature flag** in `pkg/github/github.go`:
   ```go
   if cfg.Features[FeatureMyNewFeature] {
       // Your feature code
   }
   ```

3. **Document in import-config.yaml**:
   ```yaml
   # feature_my_new_feature: Enable my new feature
   feature_my_new_feature: false
   ```

4. **Add tests**:
   ```go
   func TestMyNewFeature(t *testing.T) {
       // Test implementation
   }
   ```

## Project Structure Best Practices

- **cmd/**: CLI command definitions (Cobra)
- **pkg/github/**: GitHub API interactions
- **pkg/file/**: File I/O operations
- **pkg/compare/**: Comparison logic
- **internal/**: Private packages (if needed)

## Testing

```bash
# Unit tests
go test ./pkg/...

# Integration tests
go test ./... -tags=integration

# Coverage
go test ./... -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Benchmarks
go test -bench=. ./...
```

## CI/CD Integration

The importer is used in GitHub Actions workflows:

1. **Import workflow**: Imports single repository
2. **Bulk import workflow**: Imports multiple repositories
3. **Drift check**: Compares current vs desired state

See `.github/workflows/` in github-terraformer for workflow definitions.

## Support

- Check `ADDING_FEATURES.md` for extending functionality
- Review existing imports in `gcss-config-repo/repos/`
- Open issues in the GitHub repository
