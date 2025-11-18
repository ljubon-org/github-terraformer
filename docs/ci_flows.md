# GitHub API -> YAML -> Terraform flow

```mermaid
---
config:
  theme: forest
---
graph TD
    GH[("🌐 GitHub API<br/>(Existing Repos)")]
    GH -->|"go run main.go bulk-import"| IMP["🔧 github-repo-importer<br/>(Go CLI Tool)"]
    IMPORT_WF["📋 import.yaml workflow<br/>Single repo import"] -.-> IMP
    BULK_WF["📋 bulk-import.yaml workflow<br/>Multiple repos"] -.-> IMP
    IMP -->|"Creates YAML files"| CONFIGS["📁 configs/{owner}/*.yaml<br/>(Temporary storage)"]
    CONFIGS -->|"just import-repos"| JUST["📜 Justfile<br/>(Bash script)"]
    JUST -->|"Check if exists"| DECISION{{"File exists in<br/>repos/ ?"}}
    DECISION -->|"NO: New import"| TMP["📁 importer_tmp_dir/<br/>├── demo1.yaml<br/>└── demo2.yaml"]
    DECISION -->|"YES: Update"| REPOS["📁 repos/<br/>├── repo1.yaml<br/>└── repo2.yaml"]
    TMP -->|"local.generated_repos"| TF_IMPORT["🔄 Terraform Import Blocks<br/>import { for_each = local.generated_repos }"]
    REPOS -->|"local.new_repos"| TF_CREATE["🏗️ Terraform Resources<br/>(Create/Update)"]
    PLAN_WF["📋 tf-plan.yaml workflow<br/>On PR"] -.-> TF_IMPORT
    APPLY_WF["📋 tf-apply.yaml workflow<br/>On merge to main"] -.-> MODULE
    TF_IMPORT -->|"merge()"| MODULE["📦 module.repository<br/>for_each = local.all_repos"]
    TF_CREATE -->|"merge()"| MODULE
    MODULE -->|"terraform apply"| GH_FINAL[("🌐 GitHub<br/>(Creates/Updates/Imports)")]
    GH_FINAL -->|"After successful apply"| PROMOTE["🚀 Promote Workflow<br/>(Move files)"]
    PROMOTE_WF["📋 promote-imported-configs.yaml<br/>Manual trigger"] -.-> PROMOTE
    PROMOTE -->|"mv importer_tmp_dir/* repos/"| REPOS_FINAL["⭐ 📁 repos/<br/>ALL FILES MERGED HERE<br/>(SINGLE SOURCE OF TRUTH)"]
    REPOS -->|"Already managed repos"| REPOS_FINAL
    style GH fill:#e1f5fe
    style GH_FINAL fill:#e1f5fe
    style IMP fill:#fff3e0
    style JUST fill:#f3e5f5
    style TMP fill:#ffebee
    style REPOS fill:#e8f5e9
    style REPOS_FINAL fill:#4caf50,stroke:#2e7d32,stroke-width:3px,color:#fff
    style MODULE fill:#fce4ec
    style IMPORT_WF fill:#f0f0f0,stroke-dasharray: 5 5
    style BULK_WF fill:#f0f0f0,stroke-dasharray: 5 5
    style PLAN_WF fill:#f0f0f0,stroke-dasharray: 5 5
    style APPLY_WF fill:#f0f0f0,stroke-dasharray: 5 5
    style PROMOTE_WF fill:#f0f0f0,stroke-dasharray: 5 5
```

[MermaidLive editor](https://www.mermaidchart.com/play?utm_source=mermaid_live_editor&utm_medium=share#pako:eNqVVVGL20YQ_iuD8uILkX25i2lj2is-n2wrlS0jy4QjOowsrWT1ZK1YSXGMz5BC3lpIoIZCCfStfep7f8_9gfYndLW7kmWfAz1ssOSZ-b6Zb3Zm15KDXSS1JFmWrcjBkRf4LSsCSOdogVrgYYKS1IqY2Sd2PAfzKrcD9Ppvapb07-8_f4JekPazGbRH6jcz0rioKe-CJA0iHwwU4-TEkk5uihiQ5Ys7S_IxkCyChR1Edfo8y8JbOVjEmKSWdAfqYPQmh97-CX6QzrOZTCiQcECEk_QwdDQVTIxDyiAIaKRumNPXXRb_y0_AY-orexHCEpNbL8RLFj-m-YUIcmAomG9ArssXOQhHu5xo31ewKlkeARxkYRrEAjJ5CEYfRPEdguwUJXDdHmjgBSFK8qI7-rCr9sac7EfgvUgaa7yMENk0njJKXrqJ8iRssoIkpb8-2gkgUATTD1mSiupkkdYdvJqMTc7yGV5RhzwDjntpJ3NIHBLE6Q4xdy8SnyPnFgIPUN5gBnaldNSxqg_Xa0vqUiBhgiBikIy0Ad9Z0mbD4YoAATnUWzBES9g13xTNpxoUDZ-mi3jqBqTBMO-3n--37-kXXLTAz3e63G-3FcMZM5RlHPBeK-MWTGKXNiInNZSRXkrPkz7gyv88ypUbDrjMstUhduyw7iPaQsrkTssmmN0pP6v8oH8AExFi03FbgMqqhksae5swMq4DrPNxnCLbmcO3cBQZNmUOrKK9LCK03OPvGErbVBj_x1__-ftjJQUDJTgjDuL0NX5gG1yu3ckYae1hZTxST45DOzoyGnoEI6OciLJ0jtIejbTrfRg7jsPVcZwFIj6CFLPVUUIO9KuJpgjxC3hRPIuoneRFczdO9AcssJuFqM4kCegcrRjHA4ntMCxkuykZuHZfYuBu_Fn4pKW2rLjct9efdtVhWztYo1XJE6F50uCHYn-X8nBB0PYoBSSZQ7uWeFm44xkZ-kAXjf7tPYwIXuAUweuqsrUBfov4Lqr0lwdWehPz2GITu7LYUru5GNhRZoeQksD3ESkbJKD2cAv13j6c86diBsvRFEpZ0v1fn-BwSNuaBl1VU8YwUIyecgV9xVB4WWN12KM9GOsTo6OA3gXTmJj9k6NT0g6p5u6KnqyIrlQXylmpZMCjknRFNx29zKhgYesJeu41PbRvEr35kkN-GXCb53nn6LRqYwtXGM9R02tWjWYlEM3QHigvRlB-7TXRywfWvbReOLbXPH2WpATfotaTM_SVe34mXuVl4Kbz1nn87pmDQ0xYolU4cbpFLg56gZyDAvlFXHic5p8C3KUXjU0HYtWCJuzVJ67cR0aJTfTIqGLzPJasHIv_FSht_gNVFiuQ)

## 🔍 Detailed Process Breakdown

### 1️⃣ GitHub Import (github-repo-importer)

Location: `feature/github-repo-importer/`

```go
go run main.go bulk-import -c import-config.yaml
- Fetches repo data from GitHub API
- Creates YAML files in: configs/{owner}/*.yaml
```

### 2️⃣ YAML Storage Locations

Execution Context: Where each command runs from

```yaml
  | Stage                     | Full Path                                                                                         | Executed From Directory           |
  |---------------------------|---------------------------------------------------------------------------------------------------|-----------------------------------|
  | 1. After Import           | /home/.../github-terraformer/feature/github-repo-importer/configs/{owner}/*.yaml                  | feature/github-repo-importer/     |
  | 2. Justfile Copy - New    | /home/.../github-terraformer/feature/github-repo-provisioning/gcss_config/importer_tmp_dir/*.yaml | feature/github-repo-importer/     |
  | 3. Justfile Copy - Update | /home/.../github-terraformer/feature/github-repo-provisioning/gcss_config/repos/*.yaml            | feature/github-repo-importer/     |
  | 4. Terraform Reads        | gcss_config/importer_tmp_dir/*.yamlgcss_config/repos/*.yaml                                       | feature/github-repo-provisioning/ |
  | 5. Final (Promoted)       | /home/.../gcss-config-repo/repos/*.yaml                                                           | N/A (separate repo)               |
```

Directory Structure:

```bash
  github-terraformer/
  ├── feature/
  │   ├── github-repo-importer/           # 🔧 Import tool runs here
  │   │   ├── Justfile
  │   │   ├── main.go
  │   │   └── configs/                    # Step 1: Import creates YAMLs here
  │   │       └── {owner}/
  │   │           ├── repo1.yaml
  │   │           └── repo2.yaml
  │   │
  │   └── github-repo-provisioning/       # 📦 Terraform runs here
  │       ├── main.tf
  │       └── gcss_config/                # This is actually gcss-config-repo checkout
  │           ├── repos/                  # Step 3: Existing repos updated here
  │           │   └── existing.yaml
  │           └── importer_tmp_dir/       # Step 2: New imports placed here
  │               └── newRepo.yaml
```

Copies from: `configs/{owner}/*.yaml`

Copies to:   `../github-repo-provisioning/gcss_config/{repos or importer_tmp_dir}/`

#### Important Note: __The gcss_config/ directory is actually a checkout of the gcss-config-repo (done by GitHub Actions), not a permanent part of github-terraformer!__

### 3️⃣ YAML → Terraform Transformation

```hcl
# In main.tf - YAML becomes Terraform data
locals {
  generated_repos = {
    # Read YAML files and decode them 
    for file_path in fileset(path.module, "gcss_config/importer_tmp_dir/*.yaml") :
    basename(file_path) => yamldecode(file(file_path))  # YAML → HCL
  }
}
```

YAML structure becomes module variables

```hcl
module "repository" {
  for_each = local.all_repos

  # YAML fields map to module inputs
  name         = each.key                          # From filename
  description  = try(each.value.description, "")   # From YAML content
  visibility   = try(each.value.visibility, "")    # From YAML content
  environments = try(each.value.environments, [])  # From YAML content
}
```

### 4️⃣ Example YAML → Resource Flow

```yaml
YAML File (demo1.yaml):
description: "Demo repository for testing"
visibility: public
environments:
  - environment: production
    wait_timer: 300  # 5 minutes
    reviewers:
      users:
        - octocat
        - maintainer1
      teams:
        - platform-team
  - environment: staging
    wait_timer: 60   # 1 minute
    reviewers:
      users:
        - developer1
```

#### Becomes Terraform Resources

If in `importer_tmp_dir/` → Import block generated

```hcl
import {
    to = module.repository["demo1"].github_repository.repository
    id = "demo1"
}

import {
    to = module.repository["demo1"].github_repository_environment.environment["production"]
    id = "demo1:production"
}

import {
    to = module.repository["demo1"].github_repository_environment.environment["staging"]
    id = "demo1:staging"
}
```

Module creates actual resources

```hcl
module "repository" {
  # YAML filename → module key
  for_each = { "demo1" = <yaml_content> }

  # YAML fields → module variables
  name = "demo1"
  description = "Demo repository for testing"
  visibility = "public"
  environments = [
    {
      environment = "production"
      wait_timer = 300
      reviewers = {
        users = ["octocat", "maintainer1"]
        teams = ["platform-team"]
      }
    },
    {
      environment = "staging"
      wait_timer = 60
      reviewers = { users = ["developer1"] }
    }
  ]
}
```

### 5️⃣ Decision Tree

```yaml
  Is repo already in repos/?
  ├─ YES → Update existing file in repos/
  │        └─ Terraform updates resource
  └─ NO  → Place in importer_tmp_dir/
           ├─ Terraform imports from GitHub
           └─ After success → Move to repos/
```
