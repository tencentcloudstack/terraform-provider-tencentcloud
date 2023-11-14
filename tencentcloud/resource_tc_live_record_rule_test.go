package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveRecordRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveRecordRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_record_rule.record_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_record_rule.record_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveRecordRule = `

resource "tencentcloud_live_record_rule" "record_rule" {
  domain_name = ""
  template_id = 
  app_name = ""
  stream_name = ""
}

`
