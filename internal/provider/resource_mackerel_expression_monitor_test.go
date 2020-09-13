package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelExpressionMonitor_Basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("TerraformTestExpressionMonitor-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelExpressionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelExpressionMonitorConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "name", rName),
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
	rName := acctest.RandomWithPrefix("TerraformTestExpressionMonitor-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelExpressionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelExpressionMonitorConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "name", rName),
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
			{
				Config: testAccCheckMackerelExpressionMonitorConfigUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "name", rName),
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
	rName := acctest.RandomWithPrefix("TerraformTestExpressionMonitor-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelExpressionMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelExpressionMonitorConfigMinimum(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_expression_monitor.foobar", "name", rName),
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

func testAccCheckMackerelExpressionMonitorConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_expression_monitor" "foobar" {
    name                  = "%s"
    expression            = "avg(roleSlots(\"server:role\",\"loadavg5\"))"
    operator              = ">"
    warning               = 80.0
    critical              = 90.0
    notification_interval = 10
}`, rName)
}

func testAccCheckMackerelExpressionMonitorConfigUpdate(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_expression_monitor" "foobar" {
    name                  = "%s"
    expression            = "avg(roleSlots(\"server:role\",\"loadavg5\"))"
    operator              = ">"
    warning               = 85.5
    critical              = 95.5
    notification_interval = 10
}`, rName)
}

func testAccCheckMackerelExpressionMonitorConfigMinimum(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_expression_monitor" "foobar" {
    name                  = "%s"
    expression            = "avg(roleSlots(\"server:role\",\"loadavg5\"))"
    operator              = ">"
    warning               = 80.0
    critical              = 90.0
    notification_interval = 10
}`, rName)
}
