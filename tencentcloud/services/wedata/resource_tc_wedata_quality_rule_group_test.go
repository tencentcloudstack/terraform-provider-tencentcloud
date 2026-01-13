package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataQualityRuleGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataQualityRuleGroup,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_quality_rule_group.wedata_quality_rule_group", "id")),
		}, {
			ResourceName:      "tencentcloud_wedata_quality_rule_group.wedata_quality_rule_group",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWedataQualityRuleGroup = `

resource "tencentcloud_wedata_quality_rule_group" "wedata_quality_rule_group" {
  rule_group_exec_strategy_bo_list = {
    tasks = {
    }
    group_config = {
    }
  }
}
`
