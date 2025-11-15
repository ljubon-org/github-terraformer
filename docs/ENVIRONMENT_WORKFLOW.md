# GitHub Environment Management Workflow

This document explains how to properly manage GitHub repository environments with Terraform.

## Key Concept: Import vs Create

**IMPORTANT**: Terraform import blocks assume resources ALREADY EXIST in GitHub. You cannot "import" something that doesn't exist - you must CREATE it.

## Directory Structure and Purpose

```
gcss-config-repo/
├── importer_tmp_dir/     # For IMPORT only (existing resources)
│   └── *.yaml            # Repos that exist in GitHub
├── repos/                # For CREATE/UPDATE (ongoing management)
│   └── *.yaml           # Repos to be created or managed
```

## Workflow for Environments

### Scenario 1: Importing Existing Repository with Existing Environments

**When**: Repository and environments already exist in GitHub

**Process**:
1. Enable feature flag:
   ```yaml
   # import-config.yaml
   feature_github_environment: true
   ```

2. Run importer:
   ```bash
   go run main.go import owner/repo-name
   ```

3. File is created in `importer_tmp_dir/repo-name.yaml` with environments

4. Terraform plan will IMPORT the environments:
   ```
   Terraform will perform the following actions:
   # module.repository["repo-name"].github_repository_environment.environment["production"] will be imported
   ```

5. After successful import, move to repos/:
   ```bash
   mv gcss_config/importer_tmp_dir/repo-name.yaml gcss_config/repos/
   ```

### Scenario 2: Adding New Environments to Existing Repository

**When**: Repository exists but you want to add NEW environments

**Process**:
1. Ensure repo YAML is in `repos/` directory (NOT in `importer_tmp_dir/`)

2. Add environments to the YAML:
   ```yaml
   # repos/repo-name.yaml
   environments:
     - environment: staging
       wait_timer: 300
       can_admins_bypass: false
   ```

3. Terraform plan will CREATE the environments:
   ```
   Terraform will perform the following actions:
   # module.repository["repo-name"].github_repository_environment.environment["staging"] will be created
   ```

### Scenario 3: Creating New Repository with Environments

**When**: Both repository and environments are new

**Process**:
1. Create YAML in `repos/` directory:
   ```yaml
   # repos/new-repo.yaml
   description: "New repository"
   visibility: private
   environments:
     - environment: production
       wait_timer: 300
     - environment: staging
   ```

2. Terraform plan will CREATE both repository and environments:
   ```
   # module.repository["new-repo"].github_repository.repository will be created
   # module.repository["new-repo"].github_repository_environment.environment["production"] will be created
   # module.repository["new-repo"].github_repository_environment.environment["staging"] will be created
   ```

## Common Issues and Solutions

### Issue: "Cannot import non-existent remote object"

**Cause**: Trying to import an environment that doesn't exist in GitHub

**Solution**:
- Remove the repository YAML from `importer_tmp_dir/`
- Place it in `repos/` directory
- Terraform will CREATE instead of IMPORT

### Issue: Environments not showing in plan

**Cause**: Repository might be in both directories or in wrong directory

**Solution**:
1. Check file locations:
   ```bash
   ls gcss_config/importer_tmp_dir/test3.yaml
   ls gcss_config/repos/test3.yaml
   ```

2. Ensure it's only in ONE location:
   - `importer_tmp_dir/` = for import only
   - `repos/` = for create/manage

### Issue: Environments exist in GitHub but not in Terraform state

**Solution**: Use importer to capture them:
```bash
# Re-import the repository with environments
go run main.go import owner/repo-name

# This creates file in importer_tmp_dir/
# Terraform will import the environments
terraform apply

# Then move to repos/ for ongoing management
mv gcss_config/importer_tmp_dir/repo-name.yaml gcss_config/repos/
```

## Best Practices

1. **Never have the same repository in both directories**
   - Choose ONE: either `importer_tmp_dir/` OR `repos/`

2. **Use importer_tmp_dir/ only for initial import**
   - After import, move to `repos/`

3. **Use repos/ for all ongoing management**
   - New repositories
   - New environments
   - Configuration changes

4. **Workflow sequence**:
   ```
   Import → importer_tmp_dir/ → Apply → Move to repos/ → Manage
   ```

## Quick Decision Tree

```
Does the repository exist in GitHub?
├── NO → Create YAML in repos/ → Terraform CREATE
└── YES
    ├── Do environments exist in GitHub?
    │   ├── NO → Add to YAML in repos/ → Terraform CREATE
    │   └── YES
    │       ├── Are they in Terraform state?
    │       │   ├── NO → Import via importer_tmp_dir/ → Terraform IMPORT
    │       │   └── YES → Manage in repos/ → Terraform UPDATE
    │       └── Need to update them?
    │           └── YES → Edit YAML in repos/ → Terraform UPDATE
    └── Need to import first?
        └── YES → Use importer → importer_tmp_dir/ → Terraform IMPORT
```

## Summary

- **importer_tmp_dir/**: Temporary location for IMPORTS only
- **repos/**: Permanent location for CREATE/UPDATE/DELETE
- **Import blocks**: Only work for resources that EXIST in GitHub
- **Create**: Happens when resource is defined but doesn't exist yet
- **Never mix**: Don't have same repo in both directories