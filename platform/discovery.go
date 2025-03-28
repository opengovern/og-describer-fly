package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ConfigDiscovery represents the JSON input configuration
type ConfigDiscovery struct {
	Token string `json:"token"`
	OrganizationName string `json:"organizationName"`
}

type App struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Response struct {
	Apps []App `json:"apps"`
}

// Discover retrieves fly user info
func Discover(token string,organization_name string) ([]App, error) {
	var response Response

	url := fmt.Sprintf("https://api.machines.dev/v1/apps?org_slug=%s", organization_name)

	client := http.DefaultClient

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Apps, nil
}

func FlyIntegrationDiscovery(cfg ConfigDiscovery) ([]App, error) {
	// Check for the token
	if cfg.Token == "" {
		return nil, errors.New("token must be configured")
	}
	if cfg.OrganizationName == "" {
		return nil, errors.New("organizationName must be configured")
	}

	return Discover(cfg.Token,cfg.OrganizationName)
}
