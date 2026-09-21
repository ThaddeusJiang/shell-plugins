package bundler

import (
	"context"
	"os"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func PersonalAccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAccessToken,
		DocsURL:       sdk.URL("https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token"),
		ManagementURL: sdk.URL("https://github.com/settings/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Username,
				MarkdownDescription: "GitHub username used to authenticate to GitHub Packages.",
				Secret:              false,
			},
			{
				Name:                fieldname.Token,
				MarkdownDescription: "GitHub Personal Access Token (classic) used to authenticate to GitHub Package Registry for Bundler. Note: GitHub Packages only supports authentication using a personal access token (classic). For more information, see [Managing your personal access tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens).",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Prefix: "ghp_",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: githubPackagesProvisioner{},
		Importer:           tryGitHubPackagesEnvVar,
	}
}

const githubPackagesEnvVar = "BUNDLE_RUBYGEMS__PKG__GITHUB__COM"

type githubPackagesProvisioner struct{}

func (githubPackagesProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	out.AddEnvVar(githubPackagesEnvVar, in.ItemFields[fieldname.Username]+":"+in.ItemFields[fieldname.Token])
}

func (githubPackagesProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Environment variables are wiped automatically when the process exits.
}

func (githubPackagesProvisioner) Description() string {
	return "Provision GitHub Packages username and token as an environment variable: " + githubPackagesEnvVar
}

func tryGitHubPackagesEnvVar(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
	attempt := out.NewAttempt(importer.SourceEnvVars(githubPackagesEnvVar))
	value := os.Getenv(githubPackagesEnvVar)
	if value == "" {
		return
	}

	username, token, hasUsername := strings.Cut(value, ":")
	fields := map[sdk.FieldName]string{fieldname.Token: value}
	if hasUsername {
		if username == "" || token == "" {
			return
		}
		fields = map[sdk.FieldName]string{fieldname.Username: username, fieldname.Token: token}
	}
	// Preserve token-only imports; the required username is collected during setup.
	attempt.AddCandidate(sdk.ImportCandidate{Fields: fields})
}
