package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseCngwCanaryRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwCanaryRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_canary_rule.cngw_canary_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_canary_rule.cngw_canary_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTseCngwCanaryRule = `

resource "tencentcloud_tse_cngw_canary_rule" "cngw_canary_rule" {
  gateway_id = "gateway-xxxxxx"
  service_id = "451a9920-e67a-4519-af41-fccac0e72005"
  canary_rule {
		priority = 10
		enabled = true
		condition_list {
			type = ""
			key = ""
			operator = ""
			value = ""
			delimiter = ""
			global_config_id = ""
			global_config_name = ""
		}
		balanced_service_list {
			service_i_d = ""
			service_name = ""
			upstream_name = ""
			percent = 
		}
		service_id = ""
		service_name = ""

  }
}

`
