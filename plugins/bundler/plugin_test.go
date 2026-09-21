package bundler

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestPluginSchema(t *testing.T) {
	for _, report := range New().DeepValidate() {
		if !report.IsValid() {
			t.Errorf("invalid schema: %+v", report)
		}
	}
}

func TestBundleNeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, BundleCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"install":      {Args: []string{"install"}, ExpectedNeedsAuth: true},
		"install help": {Args: []string{"install", "--help"}},
		"help":         {Args: []string{"--help"}},
		"version":      {Args: []string{"--version"}},
		"exec":         {Args: []string{"exec", "ruby", "app.rb"}},
	})
}
