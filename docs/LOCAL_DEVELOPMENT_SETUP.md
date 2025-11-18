# Local Development Setup

This guide explains how to set up local development for the GitHub Terraformer system, which consists of two main repositories working together.

## Repository Architecture

### github-terraformer (This Repository)
Contains the tools and reusable workflows:
- **feature/github-repo-importer**: Go CLI tool for importing GitHub repositories to YAML
- **feature/github-repo-provisioning**: Terraform module for managing GitHub repositories
- **.github/workflows**: Reusable GitHub Actions workflows
- **.github/actions**: Custom GitHub Actions (gcss-config-setup, compare, pr-bot, graformer)

### gcss-config-repo (Configuration Repository)
Contains the actual configuration:
- **repos/**: Repository YAML configuration files (source of truth)
- **importer_tmp_dir/**: Temporary location for imported repositories (before Terraform import)
- **config/**: Configuration files (import-config.yaml, app-list.yaml)
- **.github/workflows**: Workflows that call github-terraformer's reusable workflows

## Prerequisites

- **Go 1.21+** (for importer development)
- **Terraform 1.0+** (for provisioning)
- **Just** command runner (`brew install just` or from https://github.com/casey/just)
- **GitHub App** or **Personal Access Token** with appropriate permissions
- **yq** for YAML processing (`brew install yq`)

## Quick Setup

### 1. Clone Both Repositories

```bash
# Clone to adjacent directories
git clone https://github.com/your-org/github-terraformer.git
git clone https://github.com/your-org/gcss-config-repo.git

# Your directory structure should be:
# workspace/
# ├── github-terraformer/
# └── gcss-config-repo/
```

### 2. Setup Configuration Files

The workflows use file copying (not symlinks) to connect the repositories:

```bash
# Copy configuration files for local development
cd github-terraformer/feature/github-repo-provisioning

# Create gcss_config directory (mimics what gcss-config-setup action does)
mkdir -p gcss_config
cp -r ../../../gcss-config-repo/* gcss_config/

# Copy required config files
cp gcss_config/config/app-list.yaml .
cp gcss_config/config/app-list.yaml ../github-repo-importer/
cp gcss_config/config/import-config.yaml ../github-repo-importer/
```

### 3. Configure Environment Variables

#### For GitHub Repo Importer (Go Tool)

```bash
cd github-terraformer/feature/github-repo-importer

# Create .env file for importer
cat > .env << 'EOF'
export GITHUB_TOKEN="ghp_your_personal_access_token"  # For local testing
export OWNER="your-org"
EOF

source .env
```

#### For Terraform Provisioning

```bash
cd github-terraformer/feature/github-repo-provisioning

# For local development with GitHub App
cat > .env << 'EOF'
# GitHub App credentials
export TF_VAR_app_id="123456"
export TF_VAR_app_installation_id="12345678"
export TF_VAR_app_private_key="$(cat ~/.secrets/github-app.pem)"

# Required variables
export TF_VAR_owner="your-org"
export TF_VAR_environment_directory="gcss_config"
EOF

source .env
```

### 4. Terraform Backend Configuration

For local development, you have two options:

#### Option A: Local State (Development)
```bash
cd github-terraformer/feature/github-repo-provisioning

# Create local backend configuration
cat > backend.tf << 'EOF'
terraform {
  backend "local" {
    path = "terraform.tfstate"
  }
}
EOF

terraform init -reconfigure
```

#### Option B: Terraform Cloud (Matches CI/CD)
```bash
cd github-terraformer/feature/github-repo-provisioning

# Use the existing backend-hcp.tf (rename if needed)
cp backend-hcp.tf backend.tf

# Set Terraform Cloud credentials
export TF_TOKEN_app_terraform_io="your-tfc-token"

terraform init
```

## Development Workflows

### Workflow 1: Importing Existing Repositories

**Step 1: Import with CLI Tool**
```bash
cd github-terraformer/feature/github-repo-importer
source .env

# Import single repository
just import-repo your-org/repo-name

# Or bulk import based on import-config.yaml
just import-repos

# Files are created in: configs/your-org/*.yaml
```

**Step 2: Copy to gcss_config**
```bash
# Copy imported files to provisioning directory
cp configs/$OWNER/*.yaml ../github-repo-provisioning/gcss_config/importer_tmp_dir/
```

**Step 3: Run Terraform Import**
```bash
cd ../github-repo-provisioning
source .env

# Review what will be imported
terraform plan

# Import the repositories
terraform apply
```

**Step 4: Promote to Permanent Location**
```bash
# After successful Terraform import, move files from importer_tmp_dir to repos
cd gcss_config
mv importer_tmp_dir/*.yaml repos/

# Commit these changes
git add -A
git commit -m "Promote imported repositories"
```

### Workflow 2: Creating New Repositories

**Step 1: Create YAML Configuration**
```bash
cd gcss-config-repo/repos

# Create new repository configuration
cat > new-repo.yaml << 'EOF'
description: "My new repository"
visibility: private
default_branch: main
has_issues: true
has_projects: false
has_wiki: false
has_downloads: true
vulnerability_alerts_enabled: true

# Optional: Add environments
environments:
  - environment: production
    wait_timer: 300
    deployment_policy:
      policy_type: protected_branches

  - environment: staging
    deployment_policy:
      policy_type: selected_branches_and_tags
      branch_patterns: ["release/*", "main"]
      tag_patterns: ["v*"]
EOF
```

**Step 2: Apply with Terraform**
```bash
cd github-terraformer/feature/github-repo-provisioning
source .env

# Update local copy
cp ../../../gcss-config-repo/repos/*.yaml gcss_config/repos/

# Plan and apply
terraform plan
terraform apply
```

## Testing

### Test GitHub Repo Importer
```bash
cd github-terraformer/feature/github-repo-importer

# Run unit tests
just test

# Test single import
GITHUB_TOKEN=$GITHUB_TOKEN go run main.go import your-org/test-repo

# Test bulk import
GITHUB_TOKEN=$GITHUB_TOKEN go run main.go bulk-import -c import-config.yaml
```

### Test Terraform Configuration
```bash
cd github-terraformer/feature/github-repo-provisioning

# Validate configuration
terraform validate

# Format check
terraform fmt -check

# Plan without applying
terraform plan
```

## CI/CD Workflow (GitHub Actions)

The actual CI/CD uses a different flow with reusable workflows:

### How It Works in CI/CD:

1. **gcss-config-repo** triggers workflows for:
   - PR creation → Bootstrap → Terraform Plan
   - Merge to main → Terraform Apply
   - Manual import → Import workflow

2. **Reusable Workflows** (in github-terraformer):
   - Called by gcss-config-repo workflows
   - Use GitHub App authentication
   - Use Terraform Cloud backend
   - Automatically handle file movements

3. **Key Actions**:
   - **gcss-config-setup**: Clones config repo and copies files
   - **compare**: Compares importer_tmp_dir with repos to identify changes
   - **pr-bot**: Creates pull requests with changes
   - **graformer**: Handles Terraform operations with HCP Terraform

### CI/CD Configuration Files

**import-config.yaml** (controls import behavior):
```yaml
# Option 1: Ignore specific repositories
ignored_repos:
  - "your-org/gcss-config-repo"  # Don't import the config repo itself
  - "your-org/github-terraformer" # Don't import the tool repo
  - "your-org/private-archived"   # Skip archived repos

# Option 2: Only import specific repositories
selected_repos:
  - "your-org/important-repo"
  - "your-org/another-repo"

# Features
feature_github_environment: true  # Enable environment import
```

**app-list.yaml** (GitHub App IDs for bypass actors):
```yaml
apps:
  - name: dependabot
    id: 29110  # GitHub's Dependabot App ID
  - name: renovate
    id: 37453  # Renovate Bot App ID
  - name: your-custom-app
    id: 123456
```

## Troubleshooting

### Authentication Issues

```bash
# Test GitHub token
curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/user

# Test GitHub App authentication (if using App)
gh api user --header "Authorization: Bearer $(gh auth token)"

# Check rate limits
curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/rate_limit
```

### File Path Issues

```bash
# Verify gcss_config structure
cd github-terraformer/feature/github-repo-provisioning
tree gcss_config -L 2

# Expected structure:
# gcss_config/
# ├── config/
# │   ├── app-list.yaml
# │   └── import-config.yaml
# ├── repos/
# │   └── *.yaml (repository configs)
# └── importer_tmp_dir/
#     └── *.yaml (pending imports)
```

### Terraform Import Failures

Common issues and solutions:

1. **"Cannot import non-existent remote object"**
   - Repository doesn't exist in GitHub yet
   - Solution: Move YAML from `importer_tmp_dir/` to `repos/` to create it

2. **"Error: Resource already exists"**
   - Repository already managed by Terraform
   - Solution: Check `terraform state list` and remove duplicate

3. **"Unauthorized"**
   - GitHub App permissions insufficient
   - Solution: Check App installation permissions

### Importer Issues

```bash
# Debug import with verbose output
cd github-terraformer/feature/github-repo-importer
go run main.go import your-org/repo -v

# Check imported file
cat configs/your-org/repo.yaml

# Validate YAML syntax
yq eval . configs/your-org/repo.yaml
```

## Local vs CI/CD Differences

| Aspect | Local Development | CI/CD (GitHub Actions) |
|--------|------------------|------------------------|
| **Authentication** | Personal token or App | GitHub App only |
| **File Management** | Manual copying | Automated via actions |
| **Terraform Backend** | Local or HCP | HCP Terraform only |
| **Config Repo** | Local directory | Checked out via action |
| **Promotion** | Manual move | Automated workflow |
| **PR Creation** | Manual | Automated via pr-bot |

## Security Best Practices

### Never Commit:
- `.env` files
- `*.pem` keys
- `terraform.tfvars`
- `*.tfstate` files
- Personal access tokens

### Secure Storage:
```bash
# Create secure directory for secrets
mkdir -p ~/.secrets && chmod 700 ~/.secrets

# Store keys securely
mv github-app.pem ~/.secrets/ && chmod 600 ~/.secrets/github-app.pem

# Use environment variables
export GITHUB_TOKEN=$(cat ~/.secrets/github-token)
```

### Use .gitignore:
```gitignore
# Add to .gitignore
.env
*.pem
terraform.tfvars
*.tfstate
*.tfstate.backup
.terraform/
```

## Next Steps

1. **Set up GitHub App**: Create an App with repository management permissions
2. **Configure Terraform Cloud**: Set up workspace for state management
3. **Test Import**: Try importing a test repository
4. **Create Repository**: Test creating a new repository via YAML
5. **Set up CI/CD**: Configure workflows in your gcss-config-repo

For more details on specific configurations, see:
- [DEVELOPERS_GUIDE.md](DEVELOPERS_GUIDE.md) - Complete YAML configuration reference
- [FEATURE_GITHUB_ENVIRONMENT.md](FEATURE_GITHUB_ENVIRONMENT.md) - Environment configuration guide
- Repository examples in `gcss-config-repo/repos/`
