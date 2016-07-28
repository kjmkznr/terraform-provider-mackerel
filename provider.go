package mackerel

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const (
	MackerelAPIKeyParamName = "MACKEREL_API_KEY"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
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
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ApiKey: d.Get("api_key").(string),
	}

	return config.NewClient()
}
