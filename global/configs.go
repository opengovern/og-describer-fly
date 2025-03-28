package global

import "github.com/opengovern/og-util/pkg/integration"

const (
	IntegrationTypeLower = "fly"                                    // example: aws, azure
	IntegrationName      = integration.Type("fly_account")          // example: AWS_ACCOUNT, AZURE_SUBSCRIPTION
	OGPluginRepoURL      = "github.com/opengovern/og-describer-fly" // example: github.com/opengovern/og-describers-aws
)

type IntegrationCredentials struct {
	Token string `json:"token"`
	OrganizationName	string `json:"organization_slug"`
}
