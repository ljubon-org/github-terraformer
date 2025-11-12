# Local Development Setup

Quick setup guide for both GitHub Repo Importer and Terraform Provisioning.

## Prerequisites

- **Go 1.21+** (for importer)
- **Terraform 1.0+** (for provisioning)
- **GitHub Token** or **GitHub App credentials**

## Quick Setup

### 1. Clone Repositories

```bash
git clone https://github.com/your-org/github-terraformer.git
git clone https://github.com/your-org/gcss-config-repo.git
```

### 2. Setup Symlinks

```bash
#!/bin/bash
# Run from workspace root

# For Terraform Provisioning
cd github-terraformer/feature/github-repo-provisioning
ln -sfn ../../../gcss-config-repo gcss_config
ln -sf gcss_config/config/app-list.yaml app-list.yaml

# For GitHub Importer
cd ../github-repo-importer
ln -sf ../../../gcss-config-repo/config/import-config.yaml import-config.yaml
ln -sf ../../../gcss-config-repo/config/app-list.yaml app-list.yaml
ln -sf ../../../gcss-config-repo/importer_tmp_dir importer_tmp_dir

echo "✅ Symlinks created"
```

### 3. Configure Environment

#### For Importer (Go Tool)

```bash
cd github-terraformer/feature/github-repo-importer

cat > .env << 'EOF'
export GITHUB_TOKEN="ghp_your_token_here"
export OWNER="your-org"
EOF

source .env
```

#### For Terraform

```bash
cd github-terraformer/feature/github-repo-provisioning

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

## Common Operations

### Import Repositories

```bash
cd github-terraformer/feature/github-repo-importer
source .env

# Single repository
just import-repo your-org/repo-name

# All repositories
just import-repos

# With custom output
go run main.go import your-org/repo -o ../custom-dir
```

### Provision with Terraform

```bash
cd github-terraformer/feature/github-repo-provisioning
source .env

# Initialize
terraform init

# Plan changes
terraform plan

# Apply changes
terraform apply

# Target specific repo
terraform apply -target='module.repository["repo-name"]'
```

## Workflow

### Creating New Repositories

1. Create YAML in `gcss-config-repo/repos/new-repo.yaml`:
```yaml
description: "My new repository"
visibility: private
default_branch: main
has_issues: true
vulnerability_alerts_enabled: true

# Optional environments
environments:
  - environment: production
    wait_timer: 300
    deployment_ref_policy:
      protected_branches_policy: true
```

2. Apply with Terraform:
```bash
cd github-terraformer/feature/github-repo-provisioning
source .env && terraform apply
```

### Importing Existing Repositories

1. Import with CLI tool:
```bash
cd github-terraformer/feature/github-repo-importer
source .env && just import-repo owner/repo
```

2. File appears in `gcss-config-repo/importer_tmp_dir/`

3. Apply with Terraform to import:
```bash
cd ../github-repo-provisioning
source .env && terraform apply
```

4. Move to permanent location:
```bash
mv gcss_config/importer_tmp_dir/repo.yaml gcss_config/repos/
```

## Configuration Files

### import-config.yaml

Controls import behavior:

```yaml
# Ignore specific repos
ignored_repos:
  - "your-org/gcss-config-repo"
  - "your-org/github-terraformer"

# Or select specific repos only
selected_repos:
  - "your-org/repo1"
  - "your-org/repo2"

# Enable environment import
feature_github_environment: true
```

### app-list.yaml

GitHub App IDs for ruleset bypass actors:

```yaml
apps:
  - name: dependabot
    id: 12345
  - name: renovate
    id: 67890
```

## Testing

```bash
# Importer tests
cd github-terraformer/feature/github-repo-importer
just test

# Terraform validation
cd ../github-repo-provisioning
terraform validate
```

## Troubleshooting

### Token/Authentication Issues

```bash
# Test GitHub token
curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/user

# Check rate limit
curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/rate_limit
```

### Symlink Issues

```bash
# Verify targets exist
ls ../../../gcss-config-repo

# Recreate if broken
rm -f symlink-name
ln -sf ../../../correct-path symlink-name
```

### Terraform Import Failures

If "Cannot import non-existent remote object":
- Repository doesn't exist in GitHub
- Move YAML from `importer_tmp_dir/` to `repos/` to create it

### Environment Variables Not Set

Always run `source .env` before commands.

## Local vs CI/CD

### Local Development

Use local Terraform state:
```bash
cp backend.tf.local backend.tf
terraform init -reconfigure
```

### CI/CD (GitHub Actions)

Uses HCP Terraform backend - workflows handle this automatically.

## Security Notes

Never commit:
- `.env` files
- `*.pem` keys
- `terraform.tfvars`
- `*.tfstate`

Store secrets in `~/.secrets/` with proper permissions:
```bash
mkdir -p ~/.secrets && chmod 700 ~/.secrets
mv *.pem ~/.secrets/ && chmod 600 ~/.secrets/*.pem
```