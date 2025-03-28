package provider

import (
	"errors"
	"github.com/opengovern/og-describer-fly/discovery/pkg/models"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"golang.org/x/net/context"
)

// DescribeListByFly A wrapper to pass Fly authorization to describers functions
func DescribeListByFly(describe func(context.Context, string, string, *models.StreamSender) ([]models.Resource, error)) models.ResourceDescriber {
	return func(ctx context.Context, cfg models.IntegrationCredentials, triggerType enums.DescribeTriggerType, additionalParameters map[string]string, stream *models.StreamSender) ([]models.Resource, error) {
		ctx = WithTriggerType(ctx, triggerType)

		var err error
		// Check for the token
		if cfg.Token == "" {
			return nil, errors.New("token must be configured")
		}
		if cfg.OrganizationName == "" {
			return nil, errors.New("organization_slug must be configured")
		}

		

		// Get values from describers
		var values []models.Resource
		result, err := describe(ctx, cfg.Token, cfg.OrganizationName, stream)
		if err != nil {
			return nil, err
		}
		values = append(values, result...)
		return values, nil
	}
}

// DescribeSingleByFly A wrapper to pass Fly authorization to describers functions
func DescribeSingleByFly(describe func(context.Context, string, string, string) (*models.Resource, error)) models.SingleResourceDescriber {
	return func(ctx context.Context, cfg models.IntegrationCredentials, triggerType enums.DescribeTriggerType, additionalParameters map[string]string, resourceID string, stream *models.StreamSender) (*models.Resource, error) {
		ctx = WithTriggerType(ctx, triggerType)

		var err error
		// Check for the token
		if cfg.Token == "" {
			return nil, errors.New("token must be configured")
		}

		

		appName := additionalParameters["AppName"]
		// Get value from describers
		value, err := describe(ctx, cfg.Token, appName, resourceID)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
}
