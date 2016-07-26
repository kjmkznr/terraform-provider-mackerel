package mackerel

import (
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/mackerelio/mackerel-client-go"
)

type Config struct {
	ApiKey string
}

func (c *Config) NewClient() (*mackerel.Client, error) {
	client := mackerel.NewClient(c.ApiKey)
	client.UserAgent = "Terraform for Mackerel"
	if logging.IsDebugOrHigher() {
		client.Verbose = true
	}
	return client, nil
}
