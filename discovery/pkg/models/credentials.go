package models

type IntegrationCredentials struct {
	Token string `json:"token"`
	OrganizationName string `json:"organization_slug"`
}
