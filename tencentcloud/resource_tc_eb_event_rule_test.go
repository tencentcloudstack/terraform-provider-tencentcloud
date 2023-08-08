package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudEbEventRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbEventRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_eb_event_rule.event_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_eb_event_rule.event_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccEbEventRule = `

resource "tencentcloud_eb_event_rule" "event_rule" {
  event_pattern = ""
  event_bus_id = ""
  rule_name = ""
  enable = 
  description = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
