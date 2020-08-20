package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"

	mackerel2 "github.com/kjmkznr/terraform-provider-mackerel"
)

func TestAccMackerelService_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mackerel2.testAccPreCheck(t) },
		Providers:    mackerel2.testAccProviders,
		CheckDestroy: testAccCheckMackerelServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelServiceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_service.foobar", "name", "foobar"),
					resource.TestCheckResourceAttr(
						"mackerel_service.foobar", "memo", "xxxxx"),
				),
			},
		},
	})
}

func testAccCheckMackerelServiceDestroy(s *terraform.State) error {
	client := mackerel2.testAccProvider.Meta().(*mackerel.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mackerel_service" {
			continue
		}

		services, err := client.FindServices()
		if err != nil {
			return fmt.Errorf("find service failed")
		}
		for _, svc := range services {
			if rs.Primary.ID == svc.Name {
				return fmt.Errorf("service still exists")
			}
		}
	}

	return nil
}

const testAccCheckMackerelServiceConfig_basic = `
resource "mackerel_service" "foobar" {
    name = "foobar"
    memo = "xxxxx"
}`
