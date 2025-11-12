# Justfile Commands Reference

This document lists all available `just` commands for the GitHub Repo Importer.

## Prerequisites

Install `just`:
```bash
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin
```

## Available Commands

### List All Commands

```bash
just --list
```

### Testing

```bash
# Run all tests
just test

# Run tests with verbose output
just test-verbose

# Run tests with coverage
just test-coverage

# Run specific test
just test-specific TestFunctionName
```

### Import Operations

```bash
# Import single repository
just import-repo owner/repository-name

# Import all repositories (bulk)
just import-repos

# Import with custom config
just import-repo-config owner/repo custom-config.yaml

# Import to custom directory
just import-repo-output owner/repo ./output-dir
```

### Comparison

```bash
# Compare two directories
just compare dir1 dir2

# Compare with verbose output
just compare-verbose dir1 dir2
```

### Build & Run

```bash
# Build binary
just build

# Run importer
just run [args...]

# Clean build artifacts
just clean
```

### Development

```bash
# Format code
just fmt

# Run linter
just lint

# Generate mocks
just mocks

# Update dependencies
just deps

# Tidy modules
just tidy
```

## Command Patterns

The Justfile typically follows these patterns:

```makefile
# Basic pattern
command-name:
    go run main.go {{args}}

# With environment variables
import-repo repo:
    #!/usr/bin/env bash
    source .env
    go run main.go import {{repo}}

# With default values
test filter='./...':
    go test {{filter}} -v

# With multiple arguments
compare dir1 dir2:
    go run main.go compare {{dir1}} {{dir2}}
```

## Creating Custom Commands

Add to `Justfile`:

```makefile
# Custom import with environments
import-with-env repo:
    #!/usr/bin/env bash
    echo "feature_github_environment: true" > temp-config.yaml
    go run main.go import {{repo}} -c temp-config.yaml
    rm temp-config.yaml

# Bulk import specific repos
import-selected repos:
    #!/usr/bin/env bash
    echo "selected_repos:" > temp-config.yaml
    for repo in {{repos}}; do
        echo "  - \"$repo\"" >> temp-config.yaml
    done
    go run main.go bulk-import -c temp-config.yaml
    rm temp-config.yaml
```

## Usage Examples

```bash
# Basic import
just import-repo myorg/myrepo

# Import multiple specific repos
just import-selected "myorg/repo1 myorg/repo2 myorg/repo3"

# Run tests for specific package
just test ./pkg/github

# Build and run
just build
./github-importer import myorg/repo

# Development workflow
just fmt           # Format code
just test          # Run tests
just lint          # Check code quality
just build         # Build binary
```

## Environment Integration

The Justfile commands automatically source `.env` when needed:

```makefile
import-repo repo:
    #!/usr/bin/env bash
    source .env  # Loads GITHUB_TOKEN and OWNER
    go run main.go import {{repo}}
```

## Tips

1. **View command before running**:
   ```bash
   just --dry-run import-repo myorg/repo
   ```

2. **Show all available commands**:
   ```bash
   just --list
   ```

3. **Get help for specific command**:
   ```bash
   just --show import-repo
   ```

4. **Run with different shell**:
   ```bash
   just --shell bash command-name
   ```

5. **Set variables**:
   ```bash
   just test filter=./pkg/github
   ```