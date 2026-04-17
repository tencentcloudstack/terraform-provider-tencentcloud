package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_config_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_config_rule.example", "config_rule_id"),
					resource.TestCheckResourceAttr("tencentcloud_config_rule.example", "identifier", "cam-user-group-bound"),
					resource.TestCheckResourceAttr("tencentcloud_config_rule.example", "risk_level", "3"),
					resource.TestCheckResourceAttr("tencentcloud_config_rule.example", "status", "ACTIVE"),
				),
			},
			{
				Config: testAccConfigRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_config_rule.example", "rule_name", "tf-example-update"),
					resource.TestCheckResourceAttr("tencentcloud_config_rule.example", "risk_level", "2"),
					resource.TestCheckResourceAttr("tencentcloud_config_rule.example", "status", "UN_ACTIVE"),
				),
			},
			{
				ResourceName:      "tencentcloud_config_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccConfigRule = `
resource "tencentcloud_config_rule" "example" {
  identifier      = "cam-user-group-bound"
  identifier_type = "SYSTEM"
  rule_name       = "tf-example"
  resource_type   = ["QCS::CAM::User"]
  risk_level      = 3
  description     = "tf example config rule"

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }

  status = "ACTIVE"
}
`

const testAccConfigRuleUpdate = `
resource "tencentcloud_config_rule" "example" {
  identifier      = "cam-user-group-bound"
  identifier_type = "SYSTEM"
  rule_name       = "tf-example-update"
  resource_type   = ["QCS::CAM::User"]
  risk_level      = 2
  description     = "updated config rule"

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }

  status = "UN_ACTIVE"
}
`
