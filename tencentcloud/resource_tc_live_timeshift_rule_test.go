package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveTimeshiftRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveTimeshiftRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_timeshift_rule.timeshift_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_timeshift_rule.timeshift_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveTimeshiftRule = `

resource "tencentcloud_live_timeshift_rule" "timeshift_rule" {
  domain_name = ""
  app_name = ""
  stream_name = ""
  template_id = 
}

`
