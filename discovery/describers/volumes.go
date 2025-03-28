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

func ListVolumes(ctx context.Context, handler *resilientbridge.ResilientBridge, org_slug string, stream *models.StreamSender) ([]models.Resource, error) {
	var wg sync.WaitGroup
	flyChan := make(chan models.Resource)
	errorChan := make(chan error, 1) // Buffered channel to capture errors

	go func() {
		defer close(flyChan)
		defer close(errorChan)
		if err := processVolumes(ctx, handler, org_slug, flyChan, &wg); err != nil {
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

func GetVolume(ctx context.Context, handler *resilientbridge.ResilientBridge, appName string, resourceID string) (*models.Resource, error) {
	volume, err := processVolume(ctx, handler, appName, resourceID)
	if err != nil {
		return nil, err
	}
	value := models.Resource{
		ID:   volume.ID,
		Name: volume.Name,
		Description: provider.VolumeDescription{
			AttachedAllocID:   volume.AttachedAllocID,
			AttachedMachineID: volume.AttachedMachineID,
			AutoBackupEnabled: volume.AutoBackupEnabled,
			BlockSize:         volume.BlockSize,
			Blocks:            volume.Blocks,
			BlocksAvail:       volume.BlocksAvail,
			BlocksFree:        volume.BlocksFree,
			CreatedAt:         volume.CreatedAt,
			Encrypted:         volume.Encrypted,
			FSType:            volume.FSType,
			HostStatus:        volume.HostStatus,
			ID:                volume.ID,
			Name:              volume.Name,
			Region:            volume.Region,
			SizeGB:            volume.SizeGB,
			SnapshotRetention: volume.SnapshotRetention,
			State:             volume.State,
			Zone:              volume.Zone,
		},
	}
	return &value, nil
}

func processVolumes(ctx context.Context, handler *resilientbridge.ResilientBridge, org_slug string, flyChan chan<- models.Resource, wg *sync.WaitGroup) error {
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
			var volumes []provider.VolumeJSON
	baseURL1 := "/apps/"

	finalURL1 := fmt.Sprintf("%s%s/volumes", baseURL1, app.Name)

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

	if err = json.Unmarshal(resp.Data, &volumes); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	for _, volume := range volumes {
		wg.Add(1)
		go func(volume provider.VolumeJSON) {
			defer wg.Done()
			value := models.Resource{
				ID:   volume.ID,
				Name: volume.Name,
				Description: provider.VolumeDescription{
					AttachedAllocID:   volume.AttachedAllocID,
					AttachedMachineID: volume.AttachedMachineID,
					AutoBackupEnabled: volume.AutoBackupEnabled,
					BlockSize:         volume.BlockSize,
					Blocks:            volume.Blocks,
					BlocksAvail:       volume.BlocksAvail,
					BlocksFree:        volume.BlocksFree,
					CreatedAt:         volume.CreatedAt,
					Encrypted:         volume.Encrypted,
					FSType:            volume.FSType,
					HostStatus:        volume.HostStatus,
					ID:                volume.ID,
					Name:              volume.Name,
					Region:            volume.Region,
					SizeGB:            volume.SizeGB,
					SnapshotRetention: volume.SnapshotRetention,
					State:             volume.State,
					Zone:              volume.Zone,
				},
			}
			flyChan <- value
		}(volume)
	}
	}
	
	

	return nil
}

func processVolume(ctx context.Context, handler *resilientbridge.ResilientBridge, appName, resourceID string) (*provider.VolumeJSON, error) {
	var volume provider.VolumeJSON
	baseURL := "/apps/"

	finalURL := fmt.Sprintf("%s%s/volumes/%s", baseURL, appName, resourceID)

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

	if err = json.Unmarshal(resp.Data, &volume); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &volume, nil
}
