package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelHostMonitor_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelHostMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMackerelHostMonitorConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "name", "terraform_for_mackerel_test_foobar"),
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
						"mackerel_host_monitor.foobar", "scopes.#", "0"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "exclude_scopes.#", "0"),
				),
			},
		},
	})
}

func TestAccMackerelHostMonitor_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelHostMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMackerelHostMonitorConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "name", "terraform_for_mackerel_test_foobar"),
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
						"mackerel_host_monitor.foobar", "scopes.#", "0"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "exclude_scopes.#", "0"),
				),
			},
			resource.TestStep{
				Config: testAccCheckMackerelHostMonitorConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "name", "terraform_for_mackerel_test_foobar_upd"),
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
						"mackerel_host_monitor.foobar", "scopes.#", "0"),
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "exclude_scopes.#", "0"),
				),
			},
		},
	})
}

func TestAccMackerelHostMonitor_Minimum(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelHostMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMackerelHostMonitorConfig_minimum,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_host_monitor.foobar", "name", "terraform_for_mackerel_test_foobar"),
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

const testAccCheckMackerelHostMonitorConfig_basic = `
resource "mackerel_host_monitor" "foobar" {
    name                  = "terraform_for_mackerel_test_foobar"
    duration              = 10
    metric                = "cpu%"
    operator              = ">"
    warning               = 80.0
    critical              = 90.0
    notification_interval = 10
}`

const testAccCheckMackerelHostMonitorConfig_update = `
resource "mackerel_host_monitor" "foobar" {
    name                  = "terraform_for_mackerel_test_foobar_upd"
    duration              = 10
    metric                = "cpu%"
    operator              = ">"
    warning               = 85.5
    critical              = 95.5
    notification_interval = 10
}`

const testAccCheckMackerelHostMonitorConfig_minimum = `
resource "mackerel_host_monitor" "foobar" {
    name                  = "terraform_for_mackerel_test_foobar"
    duration              = 10
    metric                = "cpu%"
    operator              = ">"
    warning               = 80.0
    critical              = 90.0
    notification_interval = 10
}`
