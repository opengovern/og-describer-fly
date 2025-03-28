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

func ListMachines(ctx context.Context, handler *resilientbridge.ResilientBridge, org_slug string, stream *models.StreamSender) ([]models.Resource, error) {
	var wg sync.WaitGroup
	flyChan := make(chan models.Resource)
	errorChan := make(chan error, 1) // Buffered channel to capture errors

	go func() {
		defer close(flyChan)
		defer close(errorChan)
		if err := processMachines(ctx, handler, org_slug, flyChan, &wg); err != nil {
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

func GetMachine(ctx context.Context, handler *resilientbridge.ResilientBridge, appName string, resourceID string) (*models.Resource, error) {
	machine, err := processMachine(ctx, handler, appName, resourceID)
	if err != nil {
		return nil, err
	}
	configChecks := provider.ConfigCheck{
		GracePeriod:   machine.Config.Checks.GracePeriod,
		Headers:       machine.Config.Checks.Headers,
		Interval:      machine.Config.Checks.Interval,
		Kind:          machine.Config.Checks.Kind,
		Method:        machine.Config.Checks.Method,
		Path:          machine.Config.Checks.Path,
		Port:          machine.Config.Checks.Port,
		Protocol:      machine.Config.Checks.Protocol,
		Type:          machine.Config.Checks.Type,
		Timeout:       machine.Config.Checks.Timeout,
		TLSServerName: machine.Config.Checks.TLSServerName,
		TLSSkipVerify: machine.Config.Checks.TLSSkipVerify,
	}
	var configContainers []provider.Container
	for _, container := range machine.Config.Containers {
		configContainers = append(configContainers, provider.Container{
			Cmd:        container.Cmd,
			DependsOn:  container.DependsOn,
			Env:        container.Env,
			Entrypoint: container.Entrypoint,
			Exec:       container.Exec,
			EnvFrom:    container.EnvFrom,
			Files:      container.Files,
			Image:      container.Image,
			Name:       container.Name,
			Restart:    container.Restart,
			Secrets:    container.Secrets,
			Stop:       container.Stop,
			User:       container.User,
		})
	}
	var dnsForwardRules []provider.DNSForwardRule
	for _, dnsForwardRule := range machine.Config.DNS.DNSForwardRules {
		dnsForwardRules = append(dnsForwardRules, provider.DNSForwardRule{
			Addr:     dnsForwardRule.Addr,
			Basename: dnsForwardRule.Basename,
		})
	}
	var dnsOptions []provider.DNSOption
	for _, dnsOption := range machine.Config.DNS.Options {
		dnsOptions = append(dnsOptions, provider.DNSOption{
			Name:  dnsOption.Name,
			Value: dnsOption.Value,
		})
	}
	dns := provider.DNS{
		DNSForwardRules:  dnsForwardRules,
		Hostname:         machine.Config.DNS.Hostname,
		HostnameFQDN:     machine.Config.DNS.HostnameFQDN,
		Nameservers:      machine.Config.DNS.Nameservers,
		Options:          dnsOptions,
		Searches:         machine.Config.DNS.Searches,
		SkipRegistration: machine.Config.DNS.SkipRegistration,
	}
	var files []provider.File
	for _, file := range machine.Config.Files {
		files = append(files, provider.File{
			GuestPath:  file.GuestPath,
			Mode:       file.Mode,
			RawValue:   file.RawValue,
			SecretName: file.SecretName,
		})
	}
	guest := provider.Guest{
		CPUKind:          machine.Config.Guest.CPUKind,
		CPUs:             machine.Config.Guest.CPUs,
		GPUKind:          machine.Config.Guest.GPUKind,
		GPUs:             machine.Config.Guest.GPUs,
		HostDedicationID: machine.Config.Guest.HostDedicationID,
		KernelArgs:       machine.Config.Guest.KernelArgs,
		MemoryMB:         machine.Config.Guest.MemoryMB,
	}
	init := provider.Init{
		Cmd:        machine.Config.Init.Cmd,
		Exec:       machine.Config.Init.Exec,
		Entrypoint: machine.Config.Init.Entrypoint,
		KernelArgs: machine.Config.Init.KernelArgs,
		SwapSizeMB: machine.Config.Init.SwapSizeMB,
		TTY:        machine.Config.Init.TTY,
	}
	metric := provider.Metric{
		HTTPS: machine.Config.Metrics.HTTPS,
		Port:  machine.Config.Metrics.Port,
		Path:  machine.Config.Metrics.Path,
	}
	var mounts []provider.Mount
	for _, mount := range machine.Config.Mounts {
		mounts = append(mounts, provider.Mount{
			AddSizeGB:              mount.AddSizeGB,
			Encrypted:              mount.Encrypted,
			ExtendThresholdPercent: mount.ExtendThresholdPercent,
			Name:                   mount.Name,
			Path:                   mount.Path,
			SizeGB:                 mount.SizeGB,
			SizeGBLimit:            mount.SizeGBLimit,
			Volume:                 mount.Volume,
		})
	}
	var processes []provider.Process
	for _, process := range machine.Config.Processes {
		processes = append(processes, provider.Process{
			Cmd:              process.Cmd,
			Exec:             process.Exec,
			Env:              process.Env,
			Entrypoint:       process.Entrypoint,
			EnvFrom:          process.EnvFrom,
			IgnoreAppSecrets: process.IgnoreAppSecrets,
			Secrets:          process.Secrets,
			User:             process.User,
		})
	}
	restart := provider.Restart{
		GPUBidPrice: machine.Config.Restart.GPUBidPrice,
		MaxRetries:  machine.Config.Restart.MaxRetries,
		Policy:      machine.Config.Restart.Policy,
	}
	var services []provider.Service
	for _, service := range machine.Config.Services {
		concurrency := provider.ServiceConcurrency{
			HardLimit: service.Concurrency.HardLimit,
			SoftLimit: service.Concurrency.SoftLimit,
			Type:      service.Concurrency.Type,
		}
		services = append(services, provider.Service{
			Autostart:                service.Autostart,
			Autostop:                 service.Autostop,
			Checks:                   service.Checks,
			Concurrency:              concurrency,
			ForceInstanceDescription: service.ForceInstanceDescription,
			ForceInstanceKey:         service.ForceInstanceKey,
			InternalPort:             service.InternalPort,
			MinMachinesRunning:       service.MinMachinesRunning,
			Ports:                    service.Ports,
			Protocol:                 service.Protocol,
		})
	}
	var statics []provider.Static
	for _, static := range machine.Config.Statics {
		statics = append(statics, provider.Static{
			GuestPath:     static.GuestPath,
			IndexDocument: static.IndexDocument,
			TigrisBucket:  static.TigrisBucket,
			URLPrefix:     static.URLPrefix,
		})
	}
	stopConfig := provider.StopConfig{
		Signal: machine.Config.StopConfig.Signal,
		Timeout: struct{ TimeDuration int }{
			TimeDuration: machine.Config.StopConfig.Timeout.TimeDuration,
		},
	}
	config := provider.Config{
		AutoDestroy:             machine.Config.AutoDestroy,
		Checks:                  configChecks,
		Containers:              configContainers,
		DNS:                     dns,
		DisableMachineAutostart: machine.Config.DisableMachineAutostart,
		Env:                     machine.Config.Env,
		Files:                   files,
		Guest:                   guest,
		Image:                   machine.Config.Image,
		Init:                    init,
		Metadata:                machine.Config.Metadata,
		Metrics:                 metric,
		Mounts:                  mounts,
		Processes:               processes,
		Restart:                 restart,
		Size:                    machine.Config.Size,
		Schedule:                machine.Config.Schedule,
		Services:                services,
		Standbys:                machine.Config.Standbys,
		Statics:                 statics,
		StopConfig:              stopConfig,
	}
	var checks []provider.Check
	for _, check := range machine.Checks {
		checks = append(checks, provider.Check{
			Name:      check.Name,
			Output:    check.Output,
			Status:    check.Status,
			UpdatedAt: check.UpdatedAt,
		})
	}
	var events []provider.Event
	for _, event := range machine.Events {
		events = append(events, provider.Event{
			ID:        event.ID,
			Request:   event.Request,
			Status:    event.Status,
			Source:    event.Source,
			Type:      event.Type,
			Timestamp: event.Timestamp,
		})
	}
	imageRef := provider.ImageRef{
		Digest:     machine.ImageRef.Digest,
		Labels:     machine.ImageRef.Labels,
		Registry:   machine.ImageRef.Registry,
		Repository: machine.ImageRef.Repository,
		Tag:        machine.ImageRef.Tag,
	}
	incompleteConfigChecks := provider.ConfigCheck{
		GracePeriod:   machine.IncompleteConfig.Checks.GracePeriod,
		Headers:       machine.IncompleteConfig.Checks.Headers,
		Interval:      machine.IncompleteConfig.Checks.Interval,
		Kind:          machine.IncompleteConfig.Checks.Kind,
		Method:        machine.IncompleteConfig.Checks.Method,
		Path:          machine.IncompleteConfig.Checks.Path,
		Port:          machine.IncompleteConfig.Checks.Port,
		Protocol:      machine.IncompleteConfig.Checks.Protocol,
		Type:          machine.IncompleteConfig.Checks.Type,
		Timeout:       machine.IncompleteConfig.Checks.Timeout,
		TLSServerName: machine.IncompleteConfig.Checks.TLSServerName,
		TLSSkipVerify: machine.IncompleteConfig.Checks.TLSSkipVerify,
	}
	var incompleteConfigContainers []provider.Container
	for _, container := range machine.IncompleteConfig.Containers {
		incompleteConfigContainers = append(incompleteConfigContainers, provider.Container{
			Cmd:        container.Cmd,
			DependsOn:  container.DependsOn,
			Env:        container.Env,
			Entrypoint: container.Entrypoint,
			Exec:       container.Exec,
			EnvFrom:    container.EnvFrom,
			Files:      container.Files,
			Image:      container.Image,
			Name:       container.Name,
			Restart:    container.Restart,
			Secrets:    container.Secrets,
			Stop:       container.Stop,
			User:       container.User,
		})
	}
	var incompleteDNSForwardRules []provider.DNSForwardRule
	for _, dnsForwardRule := range machine.IncompleteConfig.DNS.DNSForwardRules {
		incompleteDNSForwardRules = append(incompleteDNSForwardRules, provider.DNSForwardRule{
			Addr:     dnsForwardRule.Addr,
			Basename: dnsForwardRule.Basename,
		})
	}
	var incompleteDNSOptions []provider.DNSOption
	for _, dnsOption := range machine.IncompleteConfig.DNS.Options {
		incompleteDNSOptions = append(incompleteDNSOptions, provider.DNSOption{
			Name:  dnsOption.Name,
			Value: dnsOption.Value,
		})
	}
	incompleteDNS := provider.DNS{
		DNSForwardRules:  incompleteDNSForwardRules,
		Hostname:         machine.Config.DNS.Hostname,
		HostnameFQDN:     machine.Config.DNS.HostnameFQDN,
		Nameservers:      machine.Config.DNS.Nameservers,
		Options:          incompleteDNSOptions,
		Searches:         machine.Config.DNS.Searches,
		SkipRegistration: machine.Config.DNS.SkipRegistration,
	}
	var incompleteFiles []provider.File
	for _, file := range machine.IncompleteConfig.Files {
		incompleteFiles = append(incompleteFiles, provider.File{
			GuestPath:  file.GuestPath,
			Mode:       file.Mode,
			RawValue:   file.RawValue,
			SecretName: file.SecretName,
		})
	}
	incompleteGuest := provider.Guest{
		CPUKind:          machine.IncompleteConfig.Guest.CPUKind,
		CPUs:             machine.IncompleteConfig.Guest.CPUs,
		GPUKind:          machine.IncompleteConfig.Guest.GPUKind,
		GPUs:             machine.IncompleteConfig.Guest.GPUs,
		HostDedicationID: machine.IncompleteConfig.Guest.HostDedicationID,
		KernelArgs:       machine.IncompleteConfig.Guest.KernelArgs,
		MemoryMB:         machine.IncompleteConfig.Guest.MemoryMB,
	}
	incompleteInit := provider.Init{
		Cmd:        machine.IncompleteConfig.Init.Cmd,
		Exec:       machine.IncompleteConfig.Init.Exec,
		Entrypoint: machine.IncompleteConfig.Init.Entrypoint,
		KernelArgs: machine.IncompleteConfig.Init.KernelArgs,
		SwapSizeMB: machine.IncompleteConfig.Init.SwapSizeMB,
		TTY:        machine.IncompleteConfig.Init.TTY,
	}
	incompleteMetric := provider.Metric{
		HTTPS: machine.IncompleteConfig.Metrics.HTTPS,
		Port:  machine.IncompleteConfig.Metrics.Port,
		Path:  machine.IncompleteConfig.Metrics.Path,
	}
	var incompleteMounts []provider.Mount
	for _, mount := range machine.IncompleteConfig.Mounts {
		incompleteMounts = append(incompleteMounts, provider.Mount{
			AddSizeGB:              mount.AddSizeGB,
			Encrypted:              mount.Encrypted,
			ExtendThresholdPercent: mount.ExtendThresholdPercent,
			Name:                   mount.Name,
			Path:                   mount.Path,
			SizeGB:                 mount.SizeGB,
			SizeGBLimit:            mount.SizeGBLimit,
			Volume:                 mount.Volume,
		})
	}
	var incompleteProcesses []provider.Process
	for _, process := range machine.IncompleteConfig.Processes {
		incompleteProcesses = append(incompleteProcesses, provider.Process{
			Cmd:              process.Cmd,
			Exec:             process.Exec,
			Env:              process.Env,
			Entrypoint:       process.Entrypoint,
			EnvFrom:          process.EnvFrom,
			IgnoreAppSecrets: process.IgnoreAppSecrets,
			Secrets:          process.Secrets,
			User:             process.User,
		})
	}
	incompleteRestart := provider.Restart{
		GPUBidPrice: machine.IncompleteConfig.Restart.GPUBidPrice,
		MaxRetries:  machine.IncompleteConfig.Restart.MaxRetries,
		Policy:      machine.IncompleteConfig.Restart.Policy,
	}
	var incompleteServices []provider.Service
	for _, service := range machine.IncompleteConfig.Services {
		concurrency := provider.ServiceConcurrency{
			HardLimit: service.Concurrency.HardLimit,
			SoftLimit: service.Concurrency.SoftLimit,
			Type:      service.Concurrency.Type,
		}
		incompleteServices = append(incompleteServices, provider.Service{
			Autostart:                service.Autostart,
			Autostop:                 service.Autostop,
			Checks:                   service.Checks,
			Concurrency:              concurrency,
			ForceInstanceDescription: service.ForceInstanceDescription,
			ForceInstanceKey:         service.ForceInstanceKey,
			InternalPort:             service.InternalPort,
			MinMachinesRunning:       service.MinMachinesRunning,
			Ports:                    service.Ports,
			Protocol:                 service.Protocol,
		})
	}
	var incompleteStatics []provider.Static
	for _, static := range machine.IncompleteConfig.Statics {
		incompleteStatics = append(incompleteStatics, provider.Static{
			GuestPath:     static.GuestPath,
			IndexDocument: static.IndexDocument,
			TigrisBucket:  static.TigrisBucket,
			URLPrefix:     static.URLPrefix,
		})
	}
	incompleteStopConfig := provider.StopConfig{
		Signal: machine.IncompleteConfig.StopConfig.Signal,
		Timeout: struct{ TimeDuration int }{
			TimeDuration: machine.IncompleteConfig.StopConfig.Timeout.TimeDuration,
		},
	}
	incompleteConfig := provider.Config{
		AutoDestroy:             machine.Config.AutoDestroy,
		Checks:                  incompleteConfigChecks,
		Containers:              incompleteConfigContainers,
		DNS:                     incompleteDNS,
		DisableMachineAutostart: machine.Config.DisableMachineAutostart,
		Env:                     machine.Config.Env,
		Files:                   incompleteFiles,
		Guest:                   incompleteGuest,
		Image:                   machine.Config.Image,
		Init:                    incompleteInit,
		Metadata:                machine.Config.Metadata,
		Metrics:                 incompleteMetric,
		Mounts:                  incompleteMounts,
		Processes:               incompleteProcesses,
		Restart:                 incompleteRestart,
		Size:                    machine.Config.Size,
		Schedule:                machine.Config.Schedule,
		Services:                incompleteServices,
		Standbys:                machine.Config.Standbys,
		Statics:                 incompleteStatics,
		StopConfig:              incompleteStopConfig,
	}
	value := models.Resource{
		ID:   machine.ID,
		Name: machine.Name,
		Description: provider.MachineDescription{
			Config:           config,
			CreatedAt:        machine.CreatedAt,
			Checks:           checks,
			Events:           events,
			HostStatus:       machine.HostStatus,
			ID:               machine.ID,
			ImageRef:         imageRef,
			IncompleteConfig: incompleteConfig,
			InstanceID:       machine.InstanceID,
			Name:             machine.Name,
			Nonce:            machine.Nonce,
			PrivateIP:        machine.PrivateIP,
			Region:           machine.Region,
			State:            machine.State,
			UpdatedAt:        machine.UpdatedAt,
		},
	}
	return &value, nil
}

func processMachines(ctx context.Context, handler *resilientbridge.ResilientBridge, org_slug string, flyChan chan<- models.Resource, wg *sync.WaitGroup) error {
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
		var machines []provider.MachineJSON

	baseURL1 := "/apps/"

	finalURL1 := fmt.Sprintf("%s%s/machines", baseURL1, app.Name)

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

	if err = json.Unmarshal(resp.Data, &machines); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	for _, machine := range machines {
		wg.Add(1)
		go func(machine provider.MachineJSON) {
			defer wg.Done()
			configChecks := provider.ConfigCheck{
				GracePeriod:   machine.Config.Checks.GracePeriod,
				Headers:       machine.Config.Checks.Headers,
				Interval:      machine.Config.Checks.Interval,
				Kind:          machine.Config.Checks.Kind,
				Method:        machine.Config.Checks.Method,
				Path:          machine.Config.Checks.Path,
				Port:          machine.Config.Checks.Port,
				Protocol:      machine.Config.Checks.Protocol,
				Type:          machine.Config.Checks.Type,
				Timeout:       machine.Config.Checks.Timeout,
				TLSServerName: machine.Config.Checks.TLSServerName,
				TLSSkipVerify: machine.Config.Checks.TLSSkipVerify,
			}
			var configContainers []provider.Container
			for _, container := range machine.Config.Containers {
				configContainers = append(configContainers, provider.Container{
					Cmd:        container.Cmd,
					DependsOn:  container.DependsOn,
					Env:        container.Env,
					Entrypoint: container.Entrypoint,
					Exec:       container.Exec,
					EnvFrom:    container.EnvFrom,
					Files:      container.Files,
					Image:      container.Image,
					Name:       container.Name,
					Restart:    container.Restart,
					Secrets:    container.Secrets,
					Stop:       container.Stop,
					User:       container.User,
				})
			}
			var dnsForwardRules []provider.DNSForwardRule
			for _, dnsForwardRule := range machine.Config.DNS.DNSForwardRules {
				dnsForwardRules = append(dnsForwardRules, provider.DNSForwardRule{
					Addr:     dnsForwardRule.Addr,
					Basename: dnsForwardRule.Basename,
				})
			}
			var dnsOptions []provider.DNSOption
			for _, dnsOption := range machine.Config.DNS.Options {
				dnsOptions = append(dnsOptions, provider.DNSOption{
					Name:  dnsOption.Name,
					Value: dnsOption.Value,
				})
			}
			dns := provider.DNS{
				DNSForwardRules:  dnsForwardRules,
				Hostname:         machine.Config.DNS.Hostname,
				HostnameFQDN:     machine.Config.DNS.HostnameFQDN,
				Nameservers:      machine.Config.DNS.Nameservers,
				Options:          dnsOptions,
				Searches:         machine.Config.DNS.Searches,
				SkipRegistration: machine.Config.DNS.SkipRegistration,
			}
			var files []provider.File
			for _, file := range machine.Config.Files {
				files = append(files, provider.File{
					GuestPath:  file.GuestPath,
					Mode:       file.Mode,
					RawValue:   file.RawValue,
					SecretName: file.SecretName,
				})
			}
			guest := provider.Guest{
				CPUKind:          machine.Config.Guest.CPUKind,
				CPUs:             machine.Config.Guest.CPUs,
				GPUKind:          machine.Config.Guest.GPUKind,
				GPUs:             machine.Config.Guest.GPUs,
				HostDedicationID: machine.Config.Guest.HostDedicationID,
				KernelArgs:       machine.Config.Guest.KernelArgs,
				MemoryMB:         machine.Config.Guest.MemoryMB,
			}
			init := provider.Init{
				Cmd:        machine.Config.Init.Cmd,
				Exec:       machine.Config.Init.Exec,
				Entrypoint: machine.Config.Init.Entrypoint,
				KernelArgs: machine.Config.Init.KernelArgs,
				SwapSizeMB: machine.Config.Init.SwapSizeMB,
				TTY:        machine.Config.Init.TTY,
			}
			metric := provider.Metric{
				HTTPS: machine.Config.Metrics.HTTPS,
				Port:  machine.Config.Metrics.Port,
				Path:  machine.Config.Metrics.Path,
			}
			var mounts []provider.Mount
			for _, mount := range machine.Config.Mounts {
				mounts = append(mounts, provider.Mount{
					AddSizeGB:              mount.AddSizeGB,
					Encrypted:              mount.Encrypted,
					ExtendThresholdPercent: mount.ExtendThresholdPercent,
					Name:                   mount.Name,
					Path:                   mount.Path,
					SizeGB:                 mount.SizeGB,
					SizeGBLimit:            mount.SizeGBLimit,
					Volume:                 mount.Volume,
				})
			}
			var processes []provider.Process
			for _, process := range machine.Config.Processes {
				processes = append(processes, provider.Process{
					Cmd:              process.Cmd,
					Exec:             process.Exec,
					Env:              process.Env,
					Entrypoint:       process.Entrypoint,
					EnvFrom:          process.EnvFrom,
					IgnoreAppSecrets: process.IgnoreAppSecrets,
					Secrets:          process.Secrets,
					User:             process.User,
				})
			}
			restart := provider.Restart{
				GPUBidPrice: machine.Config.Restart.GPUBidPrice,
				MaxRetries:  machine.Config.Restart.MaxRetries,
				Policy:      machine.Config.Restart.Policy,
			}
			var services []provider.Service
			for _, service := range machine.Config.Services {
				concurrency := provider.ServiceConcurrency{
					HardLimit: service.Concurrency.HardLimit,
					SoftLimit: service.Concurrency.SoftLimit,
					Type:      service.Concurrency.Type,
				}
				services = append(services, provider.Service{
					Autostart:                service.Autostart,
					Autostop:                 service.Autostop,
					Checks:                   service.Checks,
					Concurrency:              concurrency,
					ForceInstanceDescription: service.ForceInstanceDescription,
					ForceInstanceKey:         service.ForceInstanceKey,
					InternalPort:             service.InternalPort,
					MinMachinesRunning:       service.MinMachinesRunning,
					Ports:                    service.Ports,
					Protocol:                 service.Protocol,
				})
			}
			var statics []provider.Static
			for _, static := range machine.Config.Statics {
				statics = append(statics, provider.Static{
					GuestPath:     static.GuestPath,
					IndexDocument: static.IndexDocument,
					TigrisBucket:  static.TigrisBucket,
					URLPrefix:     static.URLPrefix,
				})
			}
			stopConfig := provider.StopConfig{
				Signal: machine.Config.StopConfig.Signal,
				Timeout: struct{ TimeDuration int }{
					TimeDuration: machine.Config.StopConfig.Timeout.TimeDuration,
				},
			}
			config := provider.Config{
				AutoDestroy:             machine.Config.AutoDestroy,
				Checks:                  configChecks,
				Containers:              configContainers,
				DNS:                     dns,
				DisableMachineAutostart: machine.Config.DisableMachineAutostart,
				Env:                     machine.Config.Env,
				Files:                   files,
				Guest:                   guest,
				Image:                   machine.Config.Image,
				Init:                    init,
				Metadata:                machine.Config.Metadata,
				Metrics:                 metric,
				Mounts:                  mounts,
				Processes:               processes,
				Restart:                 restart,
				Size:                    machine.Config.Size,
				Schedule:                machine.Config.Schedule,
				Services:                services,
				Standbys:                machine.Config.Standbys,
				Statics:                 statics,
				StopConfig:              stopConfig,
			}
			var checks []provider.Check
			for _, check := range machine.Checks {
				checks = append(checks, provider.Check{
					Name:      check.Name,
					Output:    check.Output,
					Status:    check.Status,
					UpdatedAt: check.UpdatedAt,
				})
			}
			var events []provider.Event
			for _, event := range machine.Events {
				events = append(events, provider.Event{
					ID:        event.ID,
					Request:   event.Request,
					Status:    event.Status,
					Source:    event.Source,
					Type:      event.Type,
					Timestamp: event.Timestamp,
				})
			}
			imageRef := provider.ImageRef{
				Digest:     machine.ImageRef.Digest,
				Labels:     machine.ImageRef.Labels,
				Registry:   machine.ImageRef.Registry,
				Repository: machine.ImageRef.Repository,
				Tag:        machine.ImageRef.Tag,
			}
			incompleteConfigChecks := provider.ConfigCheck{
				GracePeriod:   machine.IncompleteConfig.Checks.GracePeriod,
				Headers:       machine.IncompleteConfig.Checks.Headers,
				Interval:      machine.IncompleteConfig.Checks.Interval,
				Kind:          machine.IncompleteConfig.Checks.Kind,
				Method:        machine.IncompleteConfig.Checks.Method,
				Path:          machine.IncompleteConfig.Checks.Path,
				Port:          machine.IncompleteConfig.Checks.Port,
				Protocol:      machine.IncompleteConfig.Checks.Protocol,
				Type:          machine.IncompleteConfig.Checks.Type,
				Timeout:       machine.IncompleteConfig.Checks.Timeout,
				TLSServerName: machine.IncompleteConfig.Checks.TLSServerName,
				TLSSkipVerify: machine.IncompleteConfig.Checks.TLSSkipVerify,
			}
			var incompleteConfigContainers []provider.Container
			for _, container := range machine.IncompleteConfig.Containers {
				incompleteConfigContainers = append(incompleteConfigContainers, provider.Container{
					Cmd:        container.Cmd,
					DependsOn:  container.DependsOn,
					Env:        container.Env,
					Entrypoint: container.Entrypoint,
					Exec:       container.Exec,
					EnvFrom:    container.EnvFrom,
					Files:      container.Files,
					Image:      container.Image,
					Name:       container.Name,
					Restart:    container.Restart,
					Secrets:    container.Secrets,
					Stop:       container.Stop,
					User:       container.User,
				})
			}
			var incompleteDNSForwardRules []provider.DNSForwardRule
			for _, dnsForwardRule := range machine.IncompleteConfig.DNS.DNSForwardRules {
				incompleteDNSForwardRules = append(incompleteDNSForwardRules, provider.DNSForwardRule{
					Addr:     dnsForwardRule.Addr,
					Basename: dnsForwardRule.Basename,
				})
			}
			var incompleteDNSOptions []provider.DNSOption
			for _, dnsOption := range machine.IncompleteConfig.DNS.Options {
				incompleteDNSOptions = append(incompleteDNSOptions, provider.DNSOption{
					Name:  dnsOption.Name,
					Value: dnsOption.Value,
				})
			}
			incompleteDNS := provider.DNS{
				DNSForwardRules:  incompleteDNSForwardRules,
				Hostname:         machine.Config.DNS.Hostname,
				HostnameFQDN:     machine.Config.DNS.HostnameFQDN,
				Nameservers:      machine.Config.DNS.Nameservers,
				Options:          incompleteDNSOptions,
				Searches:         machine.Config.DNS.Searches,
				SkipRegistration: machine.Config.DNS.SkipRegistration,
			}
			var incompleteFiles []provider.File
			for _, file := range machine.IncompleteConfig.Files {
				incompleteFiles = append(incompleteFiles, provider.File{
					GuestPath:  file.GuestPath,
					Mode:       file.Mode,
					RawValue:   file.RawValue,
					SecretName: file.SecretName,
				})
			}
			incompleteGuest := provider.Guest{
				CPUKind:          machine.IncompleteConfig.Guest.CPUKind,
				CPUs:             machine.IncompleteConfig.Guest.CPUs,
				GPUKind:          machine.IncompleteConfig.Guest.GPUKind,
				GPUs:             machine.IncompleteConfig.Guest.GPUs,
				HostDedicationID: machine.IncompleteConfig.Guest.HostDedicationID,
				KernelArgs:       machine.IncompleteConfig.Guest.KernelArgs,
				MemoryMB:         machine.IncompleteConfig.Guest.MemoryMB,
			}
			incompleteInit := provider.Init{
				Cmd:        machine.IncompleteConfig.Init.Cmd,
				Exec:       machine.IncompleteConfig.Init.Exec,
				Entrypoint: machine.IncompleteConfig.Init.Entrypoint,
				KernelArgs: machine.IncompleteConfig.Init.KernelArgs,
				SwapSizeMB: machine.IncompleteConfig.Init.SwapSizeMB,
				TTY:        machine.IncompleteConfig.Init.TTY,
			}
			incompleteMetric := provider.Metric{
				HTTPS: machine.IncompleteConfig.Metrics.HTTPS,
				Port:  machine.IncompleteConfig.Metrics.Port,
				Path:  machine.IncompleteConfig.Metrics.Path,
			}
			var incompleteMounts []provider.Mount
			for _, mount := range machine.IncompleteConfig.Mounts {
				incompleteMounts = append(incompleteMounts, provider.Mount{
					AddSizeGB:              mount.AddSizeGB,
					Encrypted:              mount.Encrypted,
					ExtendThresholdPercent: mount.ExtendThresholdPercent,
					Name:                   mount.Name,
					Path:                   mount.Path,
					SizeGB:                 mount.SizeGB,
					SizeGBLimit:            mount.SizeGBLimit,
					Volume:                 mount.Volume,
				})
			}
			var incompleteProcesses []provider.Process
			for _, process := range machine.IncompleteConfig.Processes {
				incompleteProcesses = append(incompleteProcesses, provider.Process{
					Cmd:              process.Cmd,
					Exec:             process.Exec,
					Env:              process.Env,
					Entrypoint:       process.Entrypoint,
					EnvFrom:          process.EnvFrom,
					IgnoreAppSecrets: process.IgnoreAppSecrets,
					Secrets:          process.Secrets,
					User:             process.User,
				})
			}
			incompleteRestart := provider.Restart{
				GPUBidPrice: machine.IncompleteConfig.Restart.GPUBidPrice,
				MaxRetries:  machine.IncompleteConfig.Restart.MaxRetries,
				Policy:      machine.IncompleteConfig.Restart.Policy,
			}
			var incompleteServices []provider.Service
			for _, service := range machine.IncompleteConfig.Services {
				concurrency := provider.ServiceConcurrency{
					HardLimit: service.Concurrency.HardLimit,
					SoftLimit: service.Concurrency.SoftLimit,
					Type:      service.Concurrency.Type,
				}
				incompleteServices = append(incompleteServices, provider.Service{
					Autostart:                service.Autostart,
					Autostop:                 service.Autostop,
					Checks:                   service.Checks,
					Concurrency:              concurrency,
					ForceInstanceDescription: service.ForceInstanceDescription,
					ForceInstanceKey:         service.ForceInstanceKey,
					InternalPort:             service.InternalPort,
					MinMachinesRunning:       service.MinMachinesRunning,
					Ports:                    service.Ports,
					Protocol:                 service.Protocol,
				})
			}
			var incompleteStatics []provider.Static
			for _, static := range machine.IncompleteConfig.Statics {
				incompleteStatics = append(incompleteStatics, provider.Static{
					GuestPath:     static.GuestPath,
					IndexDocument: static.IndexDocument,
					TigrisBucket:  static.TigrisBucket,
					URLPrefix:     static.URLPrefix,
				})
			}
			incompleteStopConfig := provider.StopConfig{
				Signal: machine.IncompleteConfig.StopConfig.Signal,
				Timeout: struct{ TimeDuration int }{
					TimeDuration: machine.IncompleteConfig.StopConfig.Timeout.TimeDuration,
				},
			}
			incompleteConfig := provider.Config{
				AutoDestroy:             machine.Config.AutoDestroy,
				Checks:                  incompleteConfigChecks,
				Containers:              incompleteConfigContainers,
				DNS:                     incompleteDNS,
				DisableMachineAutostart: machine.Config.DisableMachineAutostart,
				Env:                     machine.Config.Env,
				Files:                   incompleteFiles,
				Guest:                   incompleteGuest,
				Image:                   machine.Config.Image,
				Init:                    incompleteInit,
				Metadata:                machine.Config.Metadata,
				Metrics:                 incompleteMetric,
				Mounts:                  incompleteMounts,
				Processes:               incompleteProcesses,
				Restart:                 incompleteRestart,
				Size:                    machine.Config.Size,
				Schedule:                machine.Config.Schedule,
				Services:                incompleteServices,
				Standbys:                machine.Config.Standbys,
				Statics:                 incompleteStatics,
				StopConfig:              incompleteStopConfig,
			}
			value := models.Resource{
				ID:   machine.ID,
				Name: machine.Name,
				Description: provider.MachineDescription{
					Config:           config,
					CreatedAt:        machine.CreatedAt,
					Checks:           checks,
					Events:           events,
					HostStatus:       machine.HostStatus,
					ID:               machine.ID,
					ImageRef:         imageRef,
					IncompleteConfig: incompleteConfig,
					InstanceID:       machine.InstanceID,
					Name:             machine.Name,
					Nonce:            machine.Nonce,
					PrivateIP:        machine.PrivateIP,
					Region:           machine.Region,
					State:            machine.State,
					UpdatedAt:        machine.UpdatedAt,
				},
			}
			flyChan <- value
		}(machine)
	}

	}
	
	
	return nil
}

func processMachine(ctx context.Context, handler *resilientbridge.ResilientBridge, appName, resourceID string) (*provider.MachineJSON, error) {
	var machine provider.MachineJSON
	baseURL := "/apps/"

	finalURL := fmt.Sprintf("%s%s/machines/%s", baseURL, appName, resourceID)

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

	if err = json.Unmarshal(resp.Data, &machine); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &machine, nil
}
