package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelService_Basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("TerraformTestService-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMackerelServiceConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_service.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"mackerel_service.foobar", "memo", "xxxxx"),
				),
			},
		},
	})
}

func testAccCheckMackerelServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)

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

func testAccMackerelServiceConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foobar" {
    name = "%s"
    memo = "xxxxx"
}
`, rName)
}
