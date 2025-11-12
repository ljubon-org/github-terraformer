# Symlink Reference Guide

This document explains all symlinks needed for local development of both GitHub Repo Importer and Terraform Provisioning features.

## Directory Structure Overview

```
your-workspace/
├── github-terraformer/
│   └── feature/
│       ├── github-repo-importer/      # Go CLI tool
│       └── github-repo-provisioning/  # Terraform configuration
└── gcss-config-repo/                  # Configuration repository
    ├── repos/                         # Repository YAML files
    ├── importer_tmp_dir/              # Imported repository configs
    └── config/
        ├── app-list.yaml              # GitHub App bypass actors
        └── import-config.yaml         # Import configuration

```

## For GitHub Repo Importer

Navigate to: `github-terraformer/feature/github-repo-importer/`

### Required Symlinks

```bash
cd github-terraformer/feature/github-repo-importer

# 1. Import configuration (REQUIRED)
ln -sf ../../../gcss-config-repo/config/import-config.yaml import-config.yaml

# 2. App list configuration (REQUIRED)
ln -sf ../../../gcss-config-repo/config/app-list.yaml app-list.yaml

# 3. Output directory (OPTIONAL - for direct output to config repo)
ln -sf ../../../gcss-config-repo/importer_tmp_dir importer_tmp_dir
```

### Purpose
- **import-config.yaml**: Controls which repos to import/ignore
- **app-list.yaml**: GitHub App IDs for bypass actors in rulesets
- **importer_tmp_dir**: Where imported YAML files are written

## For Terraform Provisioning

Navigate to: `github-terraformer/feature/github-repo-provisioning/`

### Required Symlinks

```bash
cd github-terraformer/feature/github-repo-provisioning

# 1. Main configuration directory (REQUIRED)
ln -sfn ../../../gcss-config-repo gcss_config

# 2. App list configuration (REQUIRED)
ln -sf gcss_config/config/app-list.yaml app-list.yaml
```

### Purpose
- **gcss_config**: Access to all repository configurations and settings
- **app-list.yaml**: GitHub App IDs for bypass actors in rulesets

## Complete Setup Script

```bash
#!/bin/bash
# Run from your workspace root

# Setup for Importer
cd github-terraformer/feature/github-repo-importer
ln -sf ../../../gcss-config-repo/config/import-config.yaml import-config.yaml
ln -sf ../../../gcss-config-repo/config/app-list.yaml app-list.yaml
ln -sf ../../../gcss-config-repo/importer_tmp_dir importer_tmp_dir

# Setup for Terraform
cd ../github-repo-provisioning
ln -sfn ../../../gcss-config-repo gcss_config
ln -sf gcss_config/config/app-list.yaml app-list.yaml

echo "✅ All symlinks created successfully!"
```

## Verify Symlinks

```bash
# Check Importer symlinks
cd github-terraformer/feature/github-repo-importer
ls -la import-config.yaml
ls -la app-list.yaml
ls -la importer_tmp_dir

# Check Terraform symlinks
cd ../github-repo-provisioning
ls -la gcss_config
ls -la app-list.yaml
```

Expected output:
```
# Importer
lrwxrwxrwx import-config.yaml -> ../../../gcss-config-repo/config/import-config.yaml
lrwxrwxrwx app-list.yaml -> ../../../gcss-config-repo/config/app-list.yaml
lrwxrwxrwx importer_tmp_dir -> ../../../gcss-config-repo/importer_tmp_dir

# Terraform
lrwxrwxrwx gcss_config -> ../../../gcss-config-repo
lrwxrwxrwx app-list.yaml -> gcss_config/config/app-list.yaml
```

## Troubleshooting

### Broken Symlinks

If symlinks appear broken (red in `ls` output):

1. **Check relative paths**:
   ```bash
   pwd  # Verify you're in the correct directory
   ```

2. **Check target exists**:
   ```bash
   ls ../../../gcss-config-repo  # Should list the config repo contents
   ```

3. **Recreate symlinks**:
   ```bash
   rm symlink-name
   ln -sf ../../../correct-path symlink-name
   ```

### Permission Issues

```bash
# Ensure directories exist
mkdir -p ../../../gcss-config-repo/importer_tmp_dir
mkdir -p ../../../gcss-config-repo/config

# Check permissions
ls -la ../../../gcss-config-repo/
```

## Data Flow

```
1. Importer reads: import-config.yaml
   ↓
2. Importer writes: importer_tmp_dir/*.yaml
   ↓
3. Terraform reads: gcss_config/importer_tmp_dir/*.yaml
   ↓
4. Terraform imports existing repos
   ↓
5. User moves from importer_tmp_dir/ to repos/
   ↓
6. Terraform reads: gcss_config/repos/*.yaml
   ↓
7. Terraform manages repositories
```

## Important Notes

- **Never commit symlinks** that point to absolute paths
- **Use relative paths** for portability
- **The `-n` flag** in `ln -sfn` prevents following existing symlinks
- **Order matters**: Create gcss_config symlink before app-list.yaml

## Quick Reference

| Feature | Symlink | Source | Target |
|---------|---------|--------|--------|
| Importer | import-config.yaml | github-repo-importer/ | gcss-config-repo/config/import-config.yaml |
| Importer | app-list.yaml | github-repo-importer/ | gcss-config-repo/config/app-list.yaml |
| Importer | importer_tmp_dir | github-repo-importer/ | gcss-config-repo/importer_tmp_dir |
| Terraform | gcss_config | github-repo-provisioning/ | gcss-config-repo |
| Terraform | app-list.yaml | github-repo-provisioning/ | gcss_config/config/app-list.yaml |