package main

import (
	"github.com/opengovern/og-describer-fly/cloudql/fly"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: fly.Plugin})
}
