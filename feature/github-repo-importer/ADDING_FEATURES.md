# Adding New Features to the Importer

This guide shows how to add new feature-gated functionality to the importer. The architecture is designed to be extremely simple and scalable.

## Architecture Overview

The importer uses a **config-driven feature flag system**:
- Features are controlled by `config/import-config.yaml`
- Same compiled binary works for all features
- Adding features requires **zero changes to CLI or command structure**
- Perfect for CI/CD workflows

## Adding a New Feature: Step-by-Step Example

Let's add a hypothetical `feature_github_webhooks` feature that imports GitHub repository webhooks.

### Step 1: Add Feature Constant

In `pkg/github/constants.go`, add your feature constant:

```go
const (
    // Existing features
    FeatureGithubEnvironment = "feature_github_environment"

    // Your new feature
    FeatureGithubWebhooks    = "feature_github_webhooks"
)
```

### Step 2: Add Feature Logic in ImportRepo

In `pkg/github/github.go`, add your feature-gated code:

```go
func ImportRepo(repoName string, cfg *Config) (*Repository, error) {
    // ... existing code ...

    // =========================================================================
    // FEATURE: GitHub Webhooks
    // =========================================================================
    var allWebhooks []*github.Hook
    if cfg != nil && cfg.IsFeatureEnabled(FeatureGithubWebhooks) {
        webhooks, _, err := v3client.Repositories.ListHooks(
            context.Background(),
            repoNameSplit[0],
            repoNameSplit[1],
            nil,
        )
        if err != nil {
            fmt.Printf("failed to get webhooks: %v\n", err)
        } else {
            allWebhooks = webhooks

            if err := dumpManager.WriteJSONFile("webhooks.json", webhooks); err != nil {
                fmt.Printf("failed to write webhooks.json: %v\n", err)
            }
        }
    }

    // ... rest of code ...

    return &Repository{
        // ... existing fields ...
        Webhooks: resolveWebhooks(allWebhooks), // Your resolver function
    }, nil
}
```

### Step 3: Add Data Structures

In `pkg/github/repositories.go`, add the webhook field and structure:

```go
type Repository struct {
    // ... existing fields ...
    Environments []Environment         `yaml:"environments,omitempty"`
    Webhooks     []Webhook             `yaml:"webhooks,omitempty"` // New field
}

type Webhook struct {
    URL          string   `yaml:"url"`
    ContentType  string   `yaml:"content_type,omitempty"`
    Events       []string `yaml:"events,omitempty"`
    Active       bool     `yaml:"active"`
}
```

### Step 4: Document in Config File

In `gcss-config-repo/config/import-config.yaml`, document your feature:

```yaml
# =============================================================================
# FEATURE FLAGS
# =============================================================================
# Control which features the importer should use when importing repositories.
# All features are disabled by default and must be explicitly enabled.
#
# Available feature flags:
#
# feature_github_environment: Import GitHub repository environments
#   - When enabled, the importer fetches environment configurations from GitHub
#   - This includes reviewers, deployment policies, and protection rules
#   - Default: false (disabled - must explicitly enable)
#   - Usage: Set to true to enable environment import
#
# feature_github_webhooks: Import GitHub repository webhooks (example - not implemented)
#   - When enabled, the importer fetches and includes GitHub webhook configurations
#   - Webhook secrets are NOT imported (GitHub API limitation)
#   - Default: false (disabled - must explicitly enable)
#
# To enable environment import:
feature_github_environment: true

# To disable (default behavior):
#feature_github_environment: false
```

### Step 5: That's It!

**You're done!** The feature now works with:

```bash
# Enable in config/import-config.yaml
feature_github_webhooks: true

# Run importer (reads config automatically)
go run main.go import owner/repo

# OR run `Bulk import` and verify plan for all repos with webhooks which will be imported in terraform state and yaml
```

## Key Patterns

### Pattern 1: Feature Check
```go
if cfg != nil && cfg.IsFeatureEnabled(FeatureYourFeature) {
    // Your feature code
}
```

### Pattern 2: Safe Defaults
- Always check `cfg != nil` before calling methods
- Features default to `false` if not in config
- Log when features are enabled for debugging

### Pattern 3: Error Handling
- Don't fail the entire import if one feature fails
- Log errors with `fmt.Printf` for debugging
- Continue with other features

### Pattern 4: Data Persistence
- Write API responses to JSON files via `dumpManager.WriteJSONFile()`
- This helps with debugging and auditing
- Files are stored in `dumps/<repo-name>/`

## Testing Your Feature

### Test 1: Feature Disabled (Default)
```bash
# Without feature flag (should skip your feature)
go run main.go import owner/repo
# Check: dumps/owner-repo/ should NOT have your JSON file
```

### Test 2: Feature Enabled
```yaml
# In import-config.yaml
feature_your_feature: true
```

```bash
go run main.go import owner/repo
# Check: dumps/owner-repo/ should have your JSON file
# Check: YAML output should include your data
```

### Test 3: Bulk Import
```bash
go run main.go bulk-import
# Should respect feature flag for all repos
```

## Real-World Example: GitHub Environments

See the `feature_github_environment` implementation as a complete reference:

1. **Constant**: `pkg/github/constants.go:39`
   ```go
   FeatureGithubEnvironment = "feature_github_environment"
   ```

2. **Logic**: `pkg/github/github.go:105-143`
   ```go
   if cfg != nil && cfg.IsFeatureEnabled(FeatureGithubEnvironment) {
       // Fetch environments with pagination
       // Write to JSON dump
       // Store in allEnvironments
   }
   ```

3. **Data Structure**: `pkg/github/repositories.go:60-77`
   ```go
   type Environment struct {
       Environment             string
       WaitTimer               *int
       // ... more fields
   }
   ```

4. **Resolution**: `pkg/github/github.go:667-720`
   ```go
   func resolveEnvironments(envs []*github.Environment) []Environment {
       // Convert GitHub API response to YAML structure
   }
   ```

## Benefits of This Architecture

✅ **No CLI Changes** - Features are purely config-driven
✅ **Single Binary** - One build works everywhere
✅ **CI/CD Friendly** - Change config without rebuilding
✅ **Backward Compatible** - Features default to disabled
✅ **Scalable** - Add unlimited features without refactoring
✅ **Type Safe** - Constants prevent typos
✅ **Self-Documenting** - Feature names clearly describe functionality

## Common Pitfalls to Avoid

❌ **Don't add function parameters** - Use `cfg.IsFeatureEnabled()` instead
❌ **Don't add CLI flags** - Keep it config-driven
❌ **Don't fail on missing data** - Handle errors gracefully
❌ **Don't forget nil checks** - Always check `cfg != nil`
❌ **Don't skip documentation** - Update import-config.yaml

## Need Help?

- Review existing features: `FeatureGithubEnvironment`
- Check the code comments in `pkg/github/github.go:105-113`
- Look at `pkg/github/constants.go:36-43` for examples
- Read `config/import-config.yaml` for feature documentation

The architecture is intentionally simple - if you follow the 5 steps above, your feature will work perfectly with the existing system!
