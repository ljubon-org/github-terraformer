# GitHub Environment Configuration Examples

> [!IMPORTANT]
> **Environment Control via Feature Flag**
>
> GitHub Environments are controlled by a single feature flag:
>
> ### Import Control: `feature_github_environment` (Global)
> Controls whether the importer fetches environments from GitHub.
>
> ```yaml
> # gcss-config-repo/config/import-config.yaml
> feature_github_environment: true  # Required to import environments
> ```
>
> - **Location:** Import configuration file (`import-config.yaml`)
> - **Scope:** Global - affects all repository imports
> - **Default:** `false` (disabled - opt-in feature)
> - **Effect:** When `true`, environments are fetched from GitHub API and automatically managed by Terraform
> - **When to use:** Set to `true` to enable environment import and management for your organization
>
> ### How It Works
>
> | `feature_github_environment` | What Happens | Use Case |
> |------------------------------|--------------|----------|
> | `false` (or omitted) | ❌ Importer skips environments entirely<br>→ No YAML generated | Don't want environments at all |
> | `true` | ✅ Importer generates YAML with environments<br>→ Terraform manages them | Full automation (recommended) |
>
> ### Complete Workflow Example
>
> **Full Automation:**
> ```bash
> # Step 1: Enable feature in import-config.yaml
> feature_github_environment: true
>
> # Step 2: Run import
> go run main.go import myorg/my-app
>
> # Step 3: Generated repos/my-app.yaml contains:
> description: "My app"
> environments:
>   - environment: production
>     # ... environment config ...
>
> # Step 4: Run Terraform
> terraform apply
> # ✅ Terraform creates/manages the production environment
> ```
>
> **Key Point:** When `feature_github_environment` is enabled, all imported environments are automatically managed by Terraform. To exclude environments from management, simply remove them from the YAML file.

This guide shows how to configure GitHub repository environments using YAML → Terraform.

## Minimal Implementation

The simplest possible environment configuration - just creates the environment. Omit `deployment_branch_policy` to allow any branch to deploy.

**YAML Configuration** (`repos/my-app.yaml`):
```yaml
description: "My application repository"
visibility: private
default_branch: main
vulnerability_alerts_enabled: true

environments:
  - environment: staging
    # No deployment_branch_policy = any branch can deploy
```

**Generated Terraform**:
```hcl
resource "github_repository_environment" "environment" {
  environment = "staging"
  repository  = "my-app"

  # No deployment_branch_policy block = any branch can deploy
}
```

**What This Creates**:
- A simple environment named "staging"
- Any branch can deploy to it (no branch restrictions)
- No approval requirements
- No wait times

---

## Full Implementation

Complete environment configuration with all available features enabled.

**YAML Configuration** (`repos/my-app.yaml`):
```yaml
description: "My application repository"
visibility: private
default_branch: main
vulnerability_alerts_enabled: true

environments:
  - environment: production
    wait_timer: 300  # 5 minutes in seconds (300 seconds = 5 minutes)
    can_admins_bypass: false
    prevent_self_review: true
    reviewers:
      users:
        - octocat      # GitHub username
        - hubot
      teams:
        - platform-team  # Team slug
    deployment_branch_policy:
      protected_branches: true
```

**Generated Terraform**:
```hcl
# Data sources automatically resolve names to IDs
data "github_user" "reviewer" {
  for_each = toset(["octocat", "hubot"])
  username = each.value
}

data "github_team" "reviewer" {
  for_each = toset(["platform-team"])
  slug     = each.value
}

resource "github_repository_environment" "environment" {
  environment         = "production"
  repository          = "my-app"
  wait_timer          = 300
  can_admins_bypass   = false
  prevent_self_review = true

  reviewers {
    users = [for u in ["octocat", "hubot"] : data.github_user.reviewer[u].id]
    teams = [data.github_team.reviewer["platform-team"].id]
  }

  deployment_branch_policy {
    protected_branches     = true
    custom_branch_policies = false  # Always false - custom branch policies not supported
  }
}
```

**What This Creates**:
- Production environment with full protection
- 5-minute wait before deployment proceeds
- Admins cannot bypass environment protections
- Users cannot approve their own deployments
- Requires approval from one of the specified users OR teams
- Only protected branches can deploy

## Field Reference Quick Guide

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `environment` | string | Yes | Environment name | `production` |
| `wait_timer` | int | No | Wait time in seconds (max 43200) | `300` (5 min) |
| `can_admins_bypass` | bool | No | Admins bypass rules | `false` |
| `prevent_self_review` | bool | No | Prevent self-approval | `true` |
| `reviewers.users` | string[] | No | GitHub usernames (max 6) | `["octocat", "hubot"]` |
| `reviewers.teams` | string[] | No | Team slugs (max 6) | `["platform-team"]` |
| `deployment_branch_policy.protected_branches` | bool | No | Only protected branches | `true` |

**Note**: Setting `protected_branches: false` or omitting `deployment_branch_policy` entirely both allow any branch to deploy.

## Common Patterns

### Pattern 1: Least Restrictive (Dev/Test)
```yaml
environments:
  - environment: development
    can_admins_bypass: true
    # No deployment_branch_policy = any branch can deploy
```

### Pattern 2: Moderate (Staging)
```yaml
environments:
  - environment: staging
    prevent_self_review: true
    deployment_branch_policy:
      protected_branches: true
```

### Pattern 3: Most Restrictive (Production)
```yaml
environments:
  - environment: production
    wait_timer: 300  # 5 minutes
    can_admins_bypass: false
    prevent_self_review: true
    reviewers:
      users: ["octocat"]
      teams: ["platform-team"]
    deployment_branch_policy:
      protected_branches: true
```

## Troubleshooting

### Issue: "reviewers: must be 6 or fewer"
GitHub limits reviewers to 6 total (users + teams combined).

**Fix**: Reduce the number of reviewers or use teams instead of individual users.

For more details, see:
- `DEVELOPERS_GUIDE.md` - Complete field reference
- `repos/repository.yaml.example` - Template with all options
- [Terraform GitHub Provider Docs](https://registry.terraform.io/providers/integrations/github/latest/docs/resources/repository_environment)
