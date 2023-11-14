package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLivePadRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLivePadRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_pad_rule.pad_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_pad_rule.pad_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLivePadRule = `

resource "tencentcloud_live_pad_rule" "pad_rule" {
  domain_name = ""
  template_id = 
  app_name = ""
  stream_name = ""
}

`
