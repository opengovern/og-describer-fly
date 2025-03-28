package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ConfigHealthCheck represents the JSON input configuration
type ConfigHealthCheck struct {
	Token   string `json:"token"`
	AppName string `json:"app_name"`
}

// IsHealthy checks if the JWT has read access to all required resources
func IsHealthy(token string) error {
	var app App

	url := "https://api.buildpacks.io/v1/apps?org_slug=personal"

	client := http.DefaultClient

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func FlyIntegrationHealthcheck(cfg ConfigHealthCheck) (bool, error) {
	// Check for the token
	if cfg.Token == "" {
		return false, errors.New("token must be configured")
	}

	
	err := IsHealthy(cfg.Token)
	if err != nil {
		return false, err
	}

	return true, nil
}
