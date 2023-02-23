package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixChdfsLifeCycleRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsLifeCycleRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_chdfs_life_cycle_rule.life_cycle_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_chdfs_life_cycle_rule.life_cycle_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccChdfsLifeCycleRule = `

resource "tencentcloud_chdfs_life_cycle_rule" "life_cycle_rule" {
  file_system_id = "xxxx"
  life_cycle_rule {
    life_cycle_rule_name = "test"
    path                 = "/"
    transitions {
      days = 1
      type = 1
    }
    status               = 1
  }
}

`
