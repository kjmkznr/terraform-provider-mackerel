package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelHostMonitor_Basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("TerraformTestHostMonitor-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelHostMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelHostMonitorConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "notification_interval", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "max_check_attempts", "3"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "scopes.#", "0"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "exclude_scopes.#", "0"),
				),
			},
		},
	})
}

func TestAccMackerelHostMonitor_Update(t *testing.T) {
	rName := acctest.RandomWithPrefix("TerraformTestHostMonitor-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelHostMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelHostMonitorConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "notification_interval", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "max_check_attempts", "3"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "scopes.#", "0"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "exclude_scopes.#", "0"),
				),
			},
			{
				Config: testAccCheckMackerelHostMonitorConfigUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "warning", "85.5"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "critical", "95.5"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "notification_interval", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "max_check_attempts", "3"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "scopes.#", "0"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "exclude_scopes.#", "0"),
				),
			},
		},
	})
}

func TestAccMackerelHostMonitor_Minimum(t *testing.T) {
	rName := acctest.RandomWithPrefix("TerraformTestHostMonitor-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelHostMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelHostMonitorConfigMinimum(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "name", rName),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "notification_interval", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "max_check_attempts", "3"),
				),
			},
		},
	})
}

func testAccCheckMackerelHostMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mackerel_host_monitor" {
			continue
		}

		monitors, err := client.FindMonitors()
		if err != nil {
			return err
		}
		for _, monitor := range monitors {
			if monitor.MonitorID() == rs.Primary.ID {
				return fmt.Errorf("Monitor still exists")
			}
		}
	}

	return nil
}

func testAccCheckMackerelHostMonitorConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_host_monitor" "foobar" {
    name                  = "%s"
    duration              = 10
    metric                = "cpu%%"
    operator              = ">"
    warning               = 80.0
    critical              = 90.0
    notification_interval = 10
    max_check_attempts    = 3
}`, rName)
}

func testAccCheckMackerelHostMonitorConfigUpdate(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_host_monitor" "foobar" {
    name                  = "%s"
    duration              = 10
    metric                = "cpu%%"
    operator              = ">"
    warning               = 85.5
    critical              = 95.5
    notification_interval = 10
    max_check_attempts    = 3
}`, rName)
}

func testAccCheckMackerelHostMonitorConfigMinimum(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_host_monitor" "foobar" {
    name                  = "%s"
    duration              = 10
    metric                = "cpu%%"
    operator              = ">"
    warning               = 80.0
    critical              = 90.0
    notification_interval = 10
    max_check_attempts    = 3
}`, rName)
}
