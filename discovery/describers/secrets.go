package describers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/opengovern/og-describer-fly/discovery/pkg/models"
	"github.com/opengovern/og-describer-fly/discovery/provider"
	resilientbridge "github.com/opengovern/resilient-bridge"
)

func ListSecrets(ctx context.Context, handler *resilientbridge.ResilientBridge, org_slug string, stream *models.StreamSender) ([]models.Resource, error) {
	var wg sync.WaitGroup
	flyChan := make(chan models.Resource)
	errorChan := make(chan error, 1) // Buffered channel to capture errors

	go func() {
		defer close(flyChan)
		defer close(errorChan)
		if err := processSecrets(ctx, handler, org_slug, flyChan, &wg); err != nil {
			errorChan <- err // Send error to the error channel
		}
		wg.Wait()
	}()

	var values []models.Resource
	for {
		select {
		case value, ok := <-flyChan:
			if !ok {
				return values, nil
			}
			if stream != nil {
				if err := (*stream)(value); err != nil {
					return nil, err
				}
			} else {
				values = append(values, value)
			}
		case err := <-errorChan:
			return nil, err
		}
	}
}

func processSecrets(ctx context.Context, handler *resilientbridge.ResilientBridge, org_slug string, flyChan chan<- models.Resource, wg *sync.WaitGroup) error {
	var ListAppResponse provider.ListAppsResponse
	baseURL := "/apps"

	params := url.Values{}
	params.Set("org_slug", org_slug)
	finalURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req := &resilientbridge.NormalizedRequest{
		Method:   "GET",
		Endpoint: finalURL,
		Headers:  map[string]string{"Content-Type": "application/json"},
	}
	resp, err := handler.Request("fly", req)
	if err != nil {
		return fmt.Errorf("request execution failed: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("error %d: %s", resp.StatusCode, string(resp.Data))
	}

	if err = json.Unmarshal(resp.Data, &ListAppResponse); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}
	for _, app := range ListAppResponse.Apps{
		var secrets []provider.SecretJSON
	baseURL1 := "/apps/"

	finalURL1 := fmt.Sprintf("%s%s/secrets", baseURL1, app.Name)

	req := &resilientbridge.NormalizedRequest{
		Method:   "GET",
		Endpoint: finalURL1,
		Headers:  map[string]string{"accept": "application/json"},
	}

	resp, err := handler.Request("fly", req)
	if err != nil {
		return fmt.Errorf("request execution failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error %d: %s", resp.StatusCode, string(resp.Data))
	}

	if err = json.Unmarshal(resp.Data, &secrets); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	for _, secret := range secrets {
		wg.Add(1)
		go func(secret provider.SecretJSON) {
			defer wg.Done()
			value := models.Resource{
				ID:   secret.Label,
				Name: secret.Label,
				Description: provider.SecretDescription{
					Label:     secret.Label,
					PublicKey: secret.PublicKey,
					Type:      secret.Type,
				},
			}
			flyChan <- value
		}(secret)
	}

	}
	
	
	
	
	return nil
}
