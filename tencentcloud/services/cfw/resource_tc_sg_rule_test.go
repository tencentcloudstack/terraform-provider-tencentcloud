package cfw_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"testing"
)

func TestAccTencentCloudSgRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccSgRule,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sg_rule.sg_rule", "id")),
		}, {
			ResourceName:      "tencentcloud_sg_rule.sg_rule",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccSgRule = `

resource "tencentcloud_sg_rule" "sg_rule" {
  data = {
  }
}
`
