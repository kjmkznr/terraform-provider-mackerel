package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelChannelEmail_Basic(t *testing.T) {
	resourceName := "mackerel_channel.email"
	rName := acctest.RandomWithPrefix("TerraformTestChannelEmail-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMackerelChannelEmailConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
					resource.TestCheckResourceAttrSet(resourceName, "emails.0"),
					resource.TestCheckResourceAttrSet(resourceName, "events.0"),
					resource.TestCheckResourceAttrSet(resourceName, "user_ids.0"),
				),
			},
		},
	})
}

func TestAccMackerelChannelEmail_Invalid(t *testing.T) {
	rName := acctest.RandomWithPrefix("TerraformTestChannelEmail-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccMackerelChannelEmailConfigInvalid(rName),
				ExpectError: regexp.MustCompile("API request failed: invalid userIds"),
			},
		},
	})
}

func TestAccMackerelChannelSlack_Basic(t *testing.T) {
	resourceName := "mackerel_channel.slack"
	rName := acctest.RandomWithPrefix("TerraformTestChannelSlack-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMackerelChannelSlackConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "slack"),
					resource.TestCheckResourceAttrSet(resourceName, "events.0"),
					resource.TestCheckResourceAttr(resourceName, "mentions.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "mentions.ok", "ok"),
					resource.TestCheckResourceAttr(resourceName, "mentions.critical", "critical"),
					resource.TestCheckResourceAttr(resourceName, "enabled_graph_image", "true"),
				),
			},
		},
	})
}

func TestAccMackerelChannelWebhook_Basic(t *testing.T) {
	resourceName := "mackerel_channel.webhook"
	rName := acctest.RandomWithPrefix("TerraformTestChannelWebhook-")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMackerelChannelWebhookConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "webhook"),
					resource.TestCheckResourceAttrSet(resourceName, "events.0"),
					resource.TestCheckResourceAttrSet(resourceName, "events.1"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://hogehoge.com"),
					resource.TestCheckResourceAttr(resourceName, "enabled_graph_image", "true"),
				),
			},
		},
	})
}

func testAccCheckMackerelChannelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mackerel_channel" {
			continue
		}

		channels, err := client.FindChannels()
		if err != nil {
			return fmt.Errorf("find channel failed")
		}
		for _, chn := range channels {
			if rs.Primary.ID == chn.ID {
				return fmt.Errorf("channel still exists")
			}
		}
	}

	return nil
}

func testAccCheckMackerelChannelExists(resouceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*mackerel.Client)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mackerel_channel" {
				continue
			}

			channels, err := client.FindChannels()
			if err != nil {
				return fmt.Errorf("find channel failed")
			}
			for _, chn := range channels {
				if rs.Primary.ID == chn.ID {
					return nil
				}
			}
		}
		return fmt.Errorf("channel (%s) not found", resouceName)
	}
}

func testAccMackerelChannelEmailConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "email" {
  name    = "%s"
  type    = "email"
  emails  = ["foo@exapmle.com","bar@exapmle.com"]
  events  = ["alert"]
  user_ids = ["%s"]
}
`, rName, os.Getenv("USER_ID"))
}

func testAccMackerelChannelEmailConfigInvalid(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "email" {
  name    = "%s"
  type    = "email"
  emails  = ["foo@exapmle.com","bar@exapmle.com"]
  events  = ["alert"]
  user_ids = ["hoge"]
}
`, rName)
}

func testAccMackerelChannelSlackConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "slack" {
  name    = "%s"
  type    = "slack"
  events  = ["alert"]
  url = "https://hooks.slack.com/services/"
  mentions = {
    "ok": "ok",
    "critical": "critical",
  }
  enabled_graph_image = true
}
`, rName)
}

func testAccMackerelChannelWebhookConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "webhook" {
  name    = "%s"
  type    = "webhook"
  events  = ["alert", "monitor"]
  url = "https://hogehoge.com"
  enabled_graph_image = true
}
`, rName)
}
