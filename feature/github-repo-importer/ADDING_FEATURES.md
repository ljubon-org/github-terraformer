# Adding New Features to the Importer

Quick guide for adding feature-gated functionality to the importer.

## Architecture

- Features controlled by `config/import-config.yaml` under the `features:` key
- No CLI changes needed - purely config-driven
- Single binary works for all features

## Adding a Feature: 4 Steps

### Example: Adding `github_webhooks`

#### 1. Add field to `Features` struct

`pkg/github/config.go`:

```go
type Features struct {
    GithubEnvironments bool `yaml:"github_environments"`
    GithubWebhooks     bool `yaml:"github_webhooks"`  // NEW
}
```

#### 2. Add Logic

`pkg/github/github.go` in `ImportRepo()`:

```go
// =========================================================================
// FEATURE: GitHub Webhooks
// =========================================================================
var allWebhooks []*github.Hook
if cfg != nil && cfg.Features != nil && cfg.Features.GithubWebhooks {
    webhooks, _, err := v3client.Repositories.ListHooks(
        context.Background(), owner, repo, nil)
    if err != nil {
        fmt.Printf("failed to get webhooks: %v\n", err)
    } else {
        allWebhooks = webhooks
        dumpManager.WriteJSONFile("webhooks.json", webhooks)
    }
}
```

#### 3. Add Data Structure

`pkg/github/repositories.go`:

```go
type Repository struct {
    // ... existing fields ...
    Webhooks []Webhook `yaml:"webhooks,omitempty"`  // NEW
}

type Webhook struct {
    URL         string   `yaml:"url"`
    ContentType string   `yaml:"content_type,omitempty"`
    Events      []string `yaml:"events,omitempty"`
    Active      bool     `yaml:"active"`
}
```

#### 4. Document in Config

`gcss-config-repo/config/import-config.yaml` — add to the `features:` block:

```yaml
features:
  github_environments: true
  github_webhooks: true  # NEW: Import repository webhooks
```

## Usage

```bash
# Enable in import-config.yaml under features:, then:
go run main.go import owner/repo
# or
go run main.go bulk-import
```

## Key Patterns

### Always Check Config and Features

```go
if cfg != nil && cfg.Features != nil && cfg.Features.YourFeature {
    // Your code
}
```

### Error Handling

- Don't fail entire import on feature errors
- Log with `fmt.Printf`
- Continue with other features

### Data Dumps

```go
dumpManager.WriteJSONFile("feature_data.json", data)
// Saves to: dumps/<repo-name>/feature_data.json
```

## Real Example: GitHub Environments

Current implementation in `pkg/github/github.go`:

```go
if cfg != nil && cfg.Features != nil && cfg.Features.GithubEnvironments {
    // fetch environments from GitHub API...
}
```

## Testing

```bash
# 1. Feature disabled (default — features block absent or field false)
go run main.go import owner/repo
# Should NOT create dumps/owner-repo/webhooks.json

# 2. Feature enabled — in import-config.yaml:
# features:
#   github_webhooks: true
go run main.go import owner/repo
# Should create dumps/owner-repo/webhooks.json

# 3. Bulk import
go run main.go bulk-import
# Respects feature flag for all repos
```

## Do's and Don'ts

✅ **DO**

- Use direct field access on `cfg.Features`
- Guard with `cfg != nil && cfg.Features != nil`
- Handle errors gracefully
- Write dumps for debugging
- Document in `import-config.yaml` under `features:`

❌ **DON'T**

- Add string constants for feature names
- Add CLI flags or parameters
- Fail on missing data
- Skip nil checks

## Quick Reference

| File | Purpose |
|------|---------|
| `pkg/github/config.go` | Add field to `Features` struct |
| `pkg/github/github.go` | Add feature logic in `ImportRepo()` |
| `pkg/github/repositories.go` | Define data structures |
| `config/import-config.yaml` | Document & enable feature under `features:` |

That's it! Follow these 4 steps and your feature will integrate seamlessly.
