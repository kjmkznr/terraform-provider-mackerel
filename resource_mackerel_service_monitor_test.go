package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelServiceMonitor_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelServiceMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMackerelServiceMonitorConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "name", "terraform_for_mackerel_test_foobar"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "service", "Blog"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "notification_interval", "10"),
				),
			},
		},
	})
}

func TestAccMackerelServiceMonitor_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelServiceMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMackerelServiceMonitorConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "name", "terraform_for_mackerel_test_foobar"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "service", "Blog"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "notification_interval", "10"),
				),
			},
			resource.TestStep{
				Config: testAccCheckMackerelServiceMonitorConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "name", "terraform_for_mackerel_test_foobar_upd"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "service", "Blog"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "warning", "85.5"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "critical", "95.5"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "notification_interval", "10"),
				),
			},
		},
	})
}

func TestAccMackerelServiceMonitor_Minimum(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelServiceMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMackerelServiceMonitorConfig_minimum,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "name", "terraform_for_mackerel_test_foobar"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "service", "Blog"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "duration", "10"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "metric", "cpu%"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_service_monitor.foobar", "notification_interval", "10"),
				),
			},
		},
	})
}

func testAccCheckMackerelServiceMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mackerel_service_monitor" {
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

const testAccCheckMackerelServiceMonitorConfig_basic = `
resource "mackerel_service_monitor" "foobar" {
    name                  = "terraform_for_mackerel_test_foobar"
	service               = "Blog"
    duration              = 10
    metric                = "cpu%"
    operator              = ">"
    warning               = 80.0
    critical              = 90.0
    notification_interval = 10
}`

const testAccCheckMackerelServiceMonitorConfig_update = `
resource "mackerel_service_monitor" "foobar" {
    name                  = "terraform_for_mackerel_test_foobar_upd"
	service               = "Blog"
    duration              = 10
    metric                = "cpu%"
    operator              = ">"
    warning               = 85.5
    critical              = 95.5
    notification_interval = 10
}`

const testAccCheckMackerelServiceMonitorConfig_minimum = `
resource "mackerel_service_monitor" "foobar" {
    name                  = "terraform_for_mackerel_test_foobar"
	service               = "Blog"
    duration              = 10
    metric                = "cpu%"
    operator              = ">"
    warning               = 80.0
    critical              = 90.0
    notification_interval = 10
}`
