package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/kjmkznr/terraform-provider-mackerel/internal/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
