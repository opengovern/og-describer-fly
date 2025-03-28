package fly

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-fly/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableFlyApp(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name: "fly_app",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListApp,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    opengovernance.GetApp,
		},
		Columns: integrationColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description.ID"), Description: "The unique identifier for the app."},
			{Name: "machine_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Description.MachineCount"), Description: "The number of machines associated with the app."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description.Name"), Description: "The name of the app."},
			{Name: "network", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description.Network"), Description: "The network configuration of the app."},
		}),
	}
}
