package provider

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const providerName = "mackerel"

var testAccProvider *schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

var testAccProviderConfigure sync.Once

func init() {
	testAccProvider = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		providerName: func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	testAccProviderConfigure.Do(func() {
		if v := os.Getenv("MACKEREL_API_KEY"); v == "" {
			t.Fatal("MACKEREL_API_KEY must be set for acceptance tests")
		}
		err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
		if err != nil {
			t.Fatal(err)
		}
	})
}
