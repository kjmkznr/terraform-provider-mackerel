package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelExpressionMonitor_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelExpressionMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMackerelExpressionMonitorConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "name", "terraform_for_mackerel_test_foobar"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "expression", `avg(roleSlots("server:role","loadavg5"))`),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "notification_interval", "10"),
				),
			},
		},
	})
}

func TestAccMackerelExpressionMonitor_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelExpressionMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMackerelExpressionMonitorConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "name", "terraform_for_mackerel_test_foobar"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "expression", `avg(roleSlots("server:role","loadavg5"))`),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "notification_interval", "10"),
				),
			},
			resource.TestStep{
				Config: testAccCheckMackerelExpressionMonitorConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "name", "terraform_for_mackerel_test_foobar_upd"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "expression", `avg(roleSlots("server:role","loadavg5"))`),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "warning", "85.5"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "critical", "95.5"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "notification_interval", "10"),
				),
			},
		},
	})
}

func TestAccMackerelExpressionMonitor_Minimum(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelExpressionMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMackerelExpressionMonitorConfig_minimum,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "name", "terraform_for_mackerel_test_foobar"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "expression", `avg(roleSlots("server:role","loadavg5"))`),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "operator", ">"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "warning", "80"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "critical", "90"),
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "notification_interval", "10"),
				),
			},
		},
	})
}

func testAccCheckMackerelExpressionMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mackerel_expression_monitor" {
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

const testAccCheckMackerelExpressionMonitorConfig_basic = `
resource "mackerel_expression_monitor" "foobar" {
    name                  = "terraform_for_mackerel_test_foobar"
    expression            = "avg(roleSlots(\"server:role\",\"loadavg5\"))"
    operator              = ">"
    warning               = 80.0
    critical              = 90.0
    notification_interval = 10
}`

const testAccCheckMackerelExpressionMonitorConfig_update = `
resource "mackerel_expression_monitor" "foobar" {
    name                  = "terraform_for_mackerel_test_foobar_upd"
    expression            = "avg(roleSlots(\"server:role\",\"loadavg5\"))"
    operator              = ">"
    warning               = 85.5
    critical              = 95.5
    notification_interval = 10
}`

const testAccCheckMackerelExpressionMonitorConfig_minimum = `
resource "mackerel_expression_monitor" "foobar" {
    name                  = "terraform_for_mackerel_test_foobar"
    expression            = "avg(roleSlots(\"server:role\",\"loadavg5\"))"
    operator              = ">"
    warning               = 80.0
    critical              = 90.0
    notification_interval = 10
}`
