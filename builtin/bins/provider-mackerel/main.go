package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/kjmkznr/terraform-provider-mackerel"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mackerel.Provider,
	})
}
