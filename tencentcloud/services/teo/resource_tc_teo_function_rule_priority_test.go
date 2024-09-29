package teo

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoFunctionRulePriorityResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoFunctionRulePriority,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule_priority.teo_function_rule_priority", "id")),
		}, {
			ResourceName:      "tencentcloud_teo_function_rule_priority.teo_function_rule_priority",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccTeoFunctionRulePriority = `

resource "tencentcloud_teo_function_rule_priority" "teo_function_rule_priority" {
}
`
