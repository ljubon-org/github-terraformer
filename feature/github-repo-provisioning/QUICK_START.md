# Quick Start - Local Terraform Development

## One-Time Setup

```bash
# 1. Create symlinks
ln -sfn ../../../gcss-config-repo gcss_config
ln -sf gcss_config/config/app-list.yaml app-list.yaml

# 2. Create .env file (replace with your values)
cat > .env << 'EOF'
export GITHUB_TOKEN="ghp_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export GITHUB_OWNER="your-org-name"
export TF_VAR_app_private_key="$(cat ~/.secrets/github-app.pem)"
export TF_VAR_app_installation_id="12345678"
export TF_VAR_app_id="123456"
export TF_VAR_owner="your-org-name"
export TF_VAR_environment_directory="gcss_config"
EOF

# 3. Use local backend
cp backend.tf.local backend.tf

# 4. Initialize
source .env && terraform init
```

## Daily Usage

```bash
# Always start with
source .env

# Check changes
terraform plan

# Apply changes
terraform apply
```

## File Locations

- **New repos**: `gcss_config/repos/*.yaml`
- **Import existing**: `gcss_config/importer_tmp_dir/*.yaml`
- **App config**: `gcss_config/config/app-list.yaml`

## Common Fixes

```bash
# Repository doesn't exist (import error)
rm gcss_config/importer_tmp_dir/problematic-repo.yaml

# Symlink broken
ln -sfn ../../../gcss-config-repo gcss_config

# Variables not set
source .env
```
