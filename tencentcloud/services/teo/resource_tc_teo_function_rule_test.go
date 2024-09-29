package teo

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoFunctionRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoFunctionRule,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "id")),
		}, {
			ResourceName:      "tencentcloud_teo_function_rule.teo_function_rule",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccTeoFunctionRule = `

resource "tencentcloud_teo_function_rule" "teo_function_rule" {
  function_rule_conditions = ""
}
`
