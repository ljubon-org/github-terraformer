package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfigParsing_FeaturesEnabled(t *testing.T) {
	input := `
features:
  github_environments: true
`
	var cfg Config
	err := yaml.Unmarshal([]byte(input), &cfg)
	assert.NoError(t, err)
	assert.NotNil(t, cfg.Features)
	assert.True(t, cfg.Features.GithubEnvironments)
}

func TestConfigParsing_FeaturesAbsent(t *testing.T) {
	input := `
ignored_repos:
  - "org/repo"
`
	var cfg Config
	err := yaml.Unmarshal([]byte(input), &cfg)
	assert.NoError(t, err)
	assert.Nil(t, cfg.Features)
}

func TestConfigParsing_FeaturesDisabled(t *testing.T) {
	input := `
features:
  github_environments: false
`
	var cfg Config
	err := yaml.Unmarshal([]byte(input), &cfg)
	assert.NoError(t, err)
	assert.NotNil(t, cfg.Features)
	assert.False(t, cfg.Features.GithubEnvironments)
}

func TestConfigParsing_OldFlatKeyIsIgnored(t *testing.T) {
	// Old format had feature_github_environment at the root level.
	// After refactor this key is no longer captured — yaml.v3 silently ignores it.
	input := `
feature_github_environment: true
`
	var cfg Config
	err := yaml.Unmarshal([]byte(input), &cfg)
	assert.NoError(t, err)
	// Features pointer must be nil — old key is not captured by the new struct
	assert.Nil(t, cfg.Features)
}

func TestConfigParsing_UnknownFeatureKeyIsIgnored(t *testing.T) {
	// Unknown keys inside features: are silently ignored by yaml.v3
	input := `
features:
  github_environments: true
  some_future_flag: true
`
	var cfg Config
	err := yaml.Unmarshal([]byte(input), &cfg)
	assert.NoError(t, err)
	assert.NotNil(t, cfg.Features)
	assert.True(t, cfg.Features.GithubEnvironments)
}

func TestConfigValidate_BothListsError(t *testing.T) {
	cfg := Config{
		IgnoredRepos:  []string{"org/a"},
		SelectedRepos: []string{"org/b"},
	}
	assert.Error(t, cfg.Validate())
}

func TestConfigValidate_OnlyIgnoredOk(t *testing.T) {
	cfg := Config{IgnoredRepos: []string{"org/a"}}
	assert.NoError(t, cfg.Validate())
}
