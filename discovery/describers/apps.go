package describers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-describer-fly/discovery/pkg/models"
	"github.com/opengovern/og-describer-fly/discovery/provider"
	resilientbridge "github.com/opengovern/resilient-bridge"
	"net/url"
	"sync"
)

func ListApps(ctx context.Context, handler *resilientbridge.ResilientBridge,org_slug string, stream *models.StreamSender) ([]models.Resource, error) {
	var wg sync.WaitGroup
	flyChan := make(chan models.Resource)
	errorChan := make(chan error, 1) // Buffered channel to capture errors
	go func() {
		defer close(flyChan)
		defer close(errorChan)
		if err := processApps(ctx, handler, flyChan, &wg,org_slug); err != nil {
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

func GetApp(ctx context.Context, handler *resilientbridge.ResilientBridge, appName string, resourceID string) (*models.Resource, error) {
	app, err := processApp(ctx, handler, resourceID)
	if err != nil {
		return nil, err
	}
	value := models.Resource{
		ID:   app.ID,
		Name: app.Name,
		Description: provider.AppDescription{
			ID:           app.ID,
			Name:         app.Name,
			MachineCount: app.MachineCount,
			Network:      app.Network,
		},
	}
	return &value, nil
}

func processApps(ctx context.Context, handler *resilientbridge.ResilientBridge, flyChan chan<- models.Resource, wg *sync.WaitGroup,org_slug string) error {
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

	for _, app := range ListAppResponse.Apps {
		wg.Add(1)
		go func(app provider.AppJSON) {
			defer wg.Done()
			value := models.Resource{
				ID:   app.ID,
				Name: app.Name,
				Description: provider.AppDescription{
					ID:           app.ID,
					Name:         app.Name,
					MachineCount: app.MachineCount,
					Network:      app.Network,
				},
			}
			flyChan <- value
		}(app)
	}
	return nil
}

func processApp(ctx context.Context, handler *resilientbridge.ResilientBridge, resourceID string) (*provider.AppJSON, error) {
	var app provider.AppJSON
	baseURL := "/apps/"

	finalURL := fmt.Sprintf("%s%s", baseURL, resourceID)

	req := &resilientbridge.NormalizedRequest{
		Method:   "GET",
		Endpoint: finalURL,
		Headers:  map[string]string{"accept": "application/json"},
	}

	resp, err := handler.Request("fly", req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error %d: %s", resp.StatusCode, string(resp.Data))
	}

	if err = json.Unmarshal(resp.Data, &app); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &app, nil
}
