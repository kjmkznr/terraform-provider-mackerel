package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelDashboard_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelDashboardConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_dashboard.foobar", "title", "terraform_for_mackerel_test_foobar"),
					resource.TestCheckResourceAttr(
						"mackerel_dashboard.foobar", "url_path", "foo/bar"),
					resource.TestCheckResourceAttr(
						"mackerel_dashboard.foobar", "body_markdown", "# Head1\n## Head2\n\n* List1\n* List2\n"),
				),
			},
		},
	})
}

func TestAccMackerelDashboard_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelDashboardConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_dashboard.foobar", "title", "terraform_for_mackerel_test_foobar"),
					resource.TestCheckResourceAttr(
						"mackerel_dashboard.foobar", "url_path", "foo/bar"),
					resource.TestCheckResourceAttr(
						"mackerel_dashboard.foobar", "body_markdown", "# Head1\n## Head2\n\n* List1\n* List2\n"),
				),
			},
			{
				Config: testAccCheckMackerelDashboardConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"mackerel_dashboard.foobar", "title", "terraform_for_mackerel_test_foobar_upd"),
					resource.TestCheckResourceAttr(
						"mackerel_dashboard.foobar", "url_path", "bar/baz"),
					resource.TestCheckResourceAttr(
						"mackerel_dashboard.foobar", "body_markdown", "# Head1\n## Head2\n\n[Link](https://terraform.io/)\n"),
				),
			},
		},
	})
}

func testAccCheckMackerelDashboardDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mackerel_dashboard" {
			continue
		}

		_, err := client.FindDashboard(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Dashboard still exists")
		}
	}

	return nil
}

const testAccCheckMackerelDashboardConfig_basic = `
resource "mackerel_dashboard" "foobar" {
    title         = "terraform_for_mackerel_test_foobar"
    url_path      = "foo/bar"
	body_markdown = <<EOF
# Head1
## Head2

* List1
* List2
EOF
}`

const testAccCheckMackerelDashboardConfig_update = `
resource "mackerel_dashboard" "foobar" {
    title         = "terraform_for_mackerel_test_foobar_upd"
    url_path      = "bar/baz"
	body_markdown = <<EOF
# Head1
## Head2

[Link](https://terraform.io/)
EOF
}`
