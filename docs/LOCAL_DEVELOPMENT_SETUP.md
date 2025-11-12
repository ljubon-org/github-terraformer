# Local Development Setup for GitHub Terraformer

This guide explains how to set up your local environment to work with Terraform for managing GitHub repositories.

## Prerequisites

1. **Terraform** installed (v1.0+)
   ```bash
   # Ubuntu/Debian
   sudo apt update && sudo apt install terraform

   # Or download from https://www.terraform.io/downloads
   ```

2. **GitHub CLI** (optional but recommended)
   ```bash
   # Ubuntu/Debian
   sudo apt install gh
   ```

3. **GitHub App** or **Personal Access Token** with appropriate permissions:
   - Repository: Administration, Contents, Issues, Metadata, Pull requests
   - Organization: Members (read)

## Directory Structure

```
github-terraformer/
├── feature/
│   └── github-repo-provisioning/    # Main Terraform directory
│       ├── main.tf
│       ├── variables.tf
│       ├── backend.tf               # HCP Terraform backend
│       ├── backend.tf.local         # Local state backend
│       └── gcss_config/             # Symlink to config repo

gcss-config-repo/
├── repos/                           # Repository configurations
│   └── *.yaml                       # One YAML file per repository
├── importer_tmp_dir/                # Imported repository configs
└── config/
    ├── app-list.yaml                # GitHub App bypass actors
    └── import-config.yaml           # Import configuration
```

## Step 1: Create Symlinks

Navigate to the Terraform directory and create necessary symlinks:

```bash
cd github-terraformer/feature/github-repo-provisioning

# Create main symlink to config repository
ln -sfn ../../../gcss-config-repo gcss_config

# Create symlink for app-list.yaml
ln -sf gcss_config/config/app-list.yaml app-list.yaml

# Verify symlinks
ls -la | grep "^l"
```

Expected output:
```
lrwxrwxrwx 1 user user   25 Nov 14 20:25 gcss_config -> ../../../gcss-config-repo
lrwxrwxrwx 1 user user   32 Nov 14 20:33 app-list.yaml -> gcss_config/config/app-list.yaml
```

## Step 2: Configure Environment Variables

Create a `.env` file in the Terraform directory:

```bash
cat > .env << 'EOF'
# GitHub App
export TF_VAR_app_private_key="$(cat ~/.secrets/github-app-private-key.pem)"
export TF_VAR_app_installation_id="12345678"
export TF_VAR_app_id="123456"

# Required Terraform Variables
export TF_VAR_owner="your-org-name"
export TF_VAR_environment_directory="gcss_config"

# HCP Terraform Configuration (optional - for remote state)
export TF_CLOUD_ORGANIZATION="your-tfc-org"
export TF_WORKSPACE="your-workspace-name"
EOF
```

### Getting Real Values

#### For GitHub App:
1. **App ID**: Go to https://github.com/organizations/YOUR_ORG/settings/apps
2. **Installation ID**: Go to https://github.com/organizations/YOUR_ORG/settings/installations
   - Click on your app, the URL contains the installation ID
3. **Private Key**: Generate from the GitHub App settings page

## Step 3: Configure Terraform Backend

For local development, use local state:

```bash
# Use local backend for development
cp backend.tf backend.tf.hcp  # Backup HCP backend
cp backend.tf.local backend.tf

# Initialize Terraform
terraform init -reconfigure
```

Content of `backend.tf.local`:
```hcl
terraform {
  required_version = "~> 1.0"

  required_providers {
    github = {
      source = "G-Research/github"
      version = "6.5.0-gr.2"
    }
  }
}
```

## Step 4: Working with Repositories

### Understanding Repository Management

- **Files in `gcss_config/repos/`**: Repositories to be created or updated
- **Files in `gcss_config/importer_tmp_dir/`**: Repositories to be imported (must exist in GitHub)

### Creating New Repositories

1. Create a YAML file in `gcss_config/repos/`:

```yaml
# gcss_config/repos/my-new-repo.yaml
description: "My new repository"
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

# Optional: Add environments
environments:
  - environment: production
    wait_timer: 300
    can_admins_bypass: false
    prevent_self_review: true
```

2. Run Terraform:

```bash
# Load environment variables
source .env

# Check what will be created
terraform plan

# Create the repository
terraform apply
```

### Importing Existing Repositories

To import repositories that already exist in GitHub:

1. Place YAML files in `gcss_config/importer_tmp_dir/`
2. Run `terraform plan` - Terraform will import them

**Important**: If import fails with "Cannot import non-existent remote object":
- The repository doesn't exist in GitHub
- Move the YAML file to `gcss_config/repos/` to create it instead:
  ```bash
  mv gcss_config/importer_tmp_dir/repo.yaml gcss_config/repos/
  ```

## Step 5: Common Commands

```bash
# Always source environment first
source .env

# Initialize Terraform
terraform init

# Validate configuration
terraform validate

# Plan changes
terraform plan

# Apply changes
terraform apply

# Target specific repository
terraform plan -target='module.repository["repo-name"]'
terraform apply -target='module.repository["repo-name"]'

# Destroy repository (CAREFUL!)
terraform destroy -target='module.repository["repo-name"]'
```

## Troubleshooting

### Issue: "No value for required variable"

**Solution**: Source your `.env` file:
```bash
source .env
```

### Issue: "Cannot import non-existent remote object"

**Cause**: Repository doesn't exist in GitHub but file is in `importer_tmp_dir/`

**Solution**: Remove from `importer_tmp_dir/` to create instead of import:
```bash
rm gcss_config/importer_tmp_dir/problematic-repo.yaml
```

### Issue: "Function calls not allowed" in terraform.tfvars

**Cause**: Can't use `file()` function in `.tfvars` files

**Solution**: Use environment variables or heredoc syntax:
```hcl
# In terraform.tfvars (using heredoc)
app_private_key = <<-EOT
-----BEGIN RSA PRIVATE KEY-----
YOUR_KEY_CONTENT_HERE
-----END RSA PRIVATE KEY-----
EOT
```

### Issue: Symlinks not working

**Solution**: Ensure you're in the correct directory and paths are relative:
```bash
pwd  # Should be: .../github-terraformer/feature/github-repo-provisioning
ln -sfn ../../../gcss-config-repo gcss_config
```

## Environment Configuration for CI/CD

For production/CI environments, use HCP Terraform:

1. Switch to HCP backend:
   ```bash
   cp backend.tf.hcp backend.tf
   terraform init -reconfigure
   ```

2. Set workspace environment variables in HCP Terraform:
   - `app_id`
   - `app_installation_id`
   - `app_private_key`
   - `owner`
   - `environment_directory`

## Security Best Practices

1. **Never commit sensitive files**:
   - `.env`
   - `terraform.tfvars`
   - `*.pem` (private keys)
   - `backend.tf.local`

2. **Use `.gitignore`**:
   ```
   .env
   terraform.tfvars
   *.auto.tfvars
   *.pem
   backend.tf.local
   provider_override.tf
   .terraform/
   *.tfstate*
   ```

3. **Store secrets securely**:
   ```bash
   # Create secure directory
   mkdir -p ~/.secrets
   chmod 700 ~/.secrets

   # Move private keys there
   mv *.pem ~/.secrets/
   chmod 600 ~/.secrets/*.pem
   ```

## Quick Start Example

```bash
# 1. Clone repositories
git clone https://github.com/your-org/github-terraformer.git
git clone https://github.com/your-org/gcss-config-repo.git

# 2. Navigate to Terraform directory
cd github-terraformer/feature/github-repo-provisioning

# 3. Create symlinks
ln -sfn ../../../gcss-config-repo gcss_config
ln -sf gcss_config/config/app-list.yaml app-list.yaml

# 4. Create .env file with your credentials
cat > .env << 'EOF'
export GITHUB_TOKEN="ghp_your_token_here"
export GITHUB_OWNER="your-org"
export TF_VAR_owner="your-org"
export TF_VAR_environment_directory="gcss_config"
export TF_VAR_app_id="123456"
export TF_VAR_app_installation_id="12345678"
export TF_VAR_app_private_key="$(cat ~/.secrets/github-app.pem)"
EOF

# 5. Initialize and run Terraform
source .env
terraform init
terraform plan
```

## Support

For issues or questions:
- Check the [DEVELOPERS_GUIDE.md](./DEVELOPERS_GUIDE.md) for detailed configuration options
- Review repository examples in `gcss-config-repo/repos/repository.yaml.example`
- Open an issue in the GitHub repository
