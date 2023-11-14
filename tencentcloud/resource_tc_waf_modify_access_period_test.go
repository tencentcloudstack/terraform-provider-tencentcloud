package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWafModifyAccessPeriodResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafModifyAccessPeriod,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_waf_modify_access_period.modify_access_period", "id")),
			},
			{
				ResourceName:      "tencentcloud_waf_modify_access_period.modify_access_period",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafModifyAccessPeriod = `

resource "tencentcloud_waf_modify_access_period" "modify_access_period" {
  period = 
  topic_id = ""
}

`
