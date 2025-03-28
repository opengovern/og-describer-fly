package describers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/opengovern/og-describer-fly/discovery/pkg/models"
	"github.com/opengovern/og-describer-fly/discovery/provider"
)

func ListSecrets(ctx context.Context, token string, org_slug string, stream *models.StreamSender) ([]models.Resource, error) {
	var wg sync.WaitGroup
	flyChan := make(chan models.Resource)
	errorChan := make(chan error, 1) // Buffered channel to capture errors

	go func() {
		defer close(flyChan)
		defer close(errorChan)
		if err := processSecrets(ctx, token, org_slug, flyChan, &wg); err != nil {
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

func processSecrets(ctx context.Context, token string, org_slug string, flyChan chan<- models.Resource, wg *sync.WaitGroup) error {
	var ListAppResponse provider.ListAppsResponse
	baseURL := "https://api.machines.dev/v1/apps"

	params := url.Values{}
	params.Set("org_slug", org_slug)
	finalURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	// make an HTTP request
	// with http/net
	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("error %d: %s", resp.StatusCode, string(resp.Status))
	}
	body, err := io.ReadAll(resp.Body)

	if err = json.Unmarshal(body, &ListAppResponse); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}
	for _, app := range ListAppResponse.Apps{
		var secrets []provider.SecretJSON
	baseURL1 := "https://api.machines.dev/v1/apps/"

	finalURL1 := fmt.Sprintf("%s%s/secrets", baseURL1, app.Name)
	req, err := http.NewRequest("GET", finalURL1, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("error %d: %s", resp.StatusCode, string(resp.Status))
	}
	body, err := io.ReadAll(resp.Body)



	if err != nil {
		return fmt.Errorf("request execution failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error %d: %s", resp.StatusCode, string(body))
	}

	if err = json.Unmarshal(body, &secrets); err != nil {
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
