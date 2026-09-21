package bundler

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:    "ghp_1234567890123456789012345678901234567890",
				fieldname.Username: "octocat",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"BUNDLE_RUBYGEMS__PKG__GITHUB__COM": "octocat:ghp_1234567890123456789012345678901234567890",
				},
			},
		},
	})
}

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"empty": {},
		"missing token": {
			Environment: map[string]string{"BUNDLE_RUBYGEMS__PKG__GITHUB__COM": "octocat:"},
		},
		"missing username": {
			Environment: map[string]string{"BUNDLE_RUBYGEMS__PKG__GITHUB__COM": ":ghp_test"},
		},
		"token only": {
			Environment: map[string]string{"BUNDLE_RUBYGEMS__PKG__GITHUB__COM": "ghp_test"},
			ExpectedCandidates: []sdk.ImportCandidate{{
				Fields: map[sdk.FieldName]string{fieldname.Token: "ghp_test"},
			}},
		},
		"environment": {
			Environment: map[string]string{
				"BUNDLE_RUBYGEMS__PKG__GITHUB__COM": "octocat:ghp_1234567890123456789012345678901234567890",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:    "ghp_1234567890123456789012345678901234567890",
						fieldname.Username: "octocat",
					},
				},
			},
		},
	})
}
