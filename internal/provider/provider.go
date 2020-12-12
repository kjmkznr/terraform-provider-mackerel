package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	MackerelAPIKeyParamName = "MACKEREL_API_KEY"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(MackerelAPIKeyParamName, nil),
				Description: "your Mackerel APIKey",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mackerel_host_monitor":       resourceMackerelHostMonitor(),
			"mackerel_service_monitor":    resourceMackerelServiceMonitor(),
			"mackerel_external_monitor":   resourceMackerelExternalMonitor(),
			"mackerel_expression_monitor": resourceMackerelExpressionMonitor(),
			"mackerel_dashboard":          resourceMackerelDashboard(),
			"mackerel_service":            resourceMackerelService(),
			"mackerel_channel":            resourceMackerelChannel(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	config := Config{
		ApiKey: d.Get("api_key").(string),
	}

	c, err := config.NewClient()
	if err != nil {
		diag.FromErr(err)
	}
	return c, diags
}
