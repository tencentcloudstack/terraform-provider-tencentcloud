package tse_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -test.run TestAccTencentCloudTseCngwStrategyBindGroupResource_basic -v -timeout=0
func TestAccTencentCloudTseCngwStrategyBindGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwStrategyBindGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_strategy_bind_group.cngw_strategy_bind_group", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_strategy_bind_group.cngw_strategy_bind_group", "strategy_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_strategy_bind_group.cngw_strategy_bind_group", "group_id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy_bind_group.cngw_strategy_bind_group", "option", "bind"),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_strategy_bind_group.cngw_strategy_bind_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTseCngwStrategyBindGroupUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_strategy_bind_group.cngw_strategy_bind_group", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_strategy_bind_group.cngw_strategy_bind_group", "strategy_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_strategy_bind_group.cngw_strategy_bind_group", "group_id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_strategy_bind_group.cngw_strategy_bind_group", "option", "unbind"),
				),
			},
		},
	})
}

const testAccTseCngwBindGroup = `
resource "tencentcloud_tse_cngw_group" "cngw_group" {
	description = "terraform desc"
	gateway_id  = tencentcloud_tse_cngw_gateway.cngw_gateway.id
	name        = "terraform-group"
	subnet_id   = tencentcloud_subnet.subnet.id
  
	node_config {
	  number        = 2
	  specification = "1c2g"
	}
  }
`

const testAccTseCngwStrategyBindGroup = testAccTseCngwStrategy + testAccTseCngwBindGroup + `

resource "tencentcloud_tse_cngw_strategy_bind_group" "cngw_strategy_bind_group" {
  gateway_id = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  strategy_id = tencentcloud_tse_cngw_strategy.cngw_strategy.strategy_id
  group_id = tencentcloud_tse_cngw_group.cngw_group.group_id
  option      = "bind"
}

`

const testAccTseCngwStrategyBindGroupUp = testAccTseCngwStrategy + testAccTseCngwBindGroup + `

resource "tencentcloud_tse_cngw_strategy_bind_group" "cngw_strategy_bind_group" {
  gateway_id = tencentcloud_tse_cngw_gateway.cngw_gateway.id
  strategy_id = tencentcloud_tse_cngw_strategy.cngw_strategy.strategy_id
  group_id = tencentcloud_tse_cngw_group.cngw_group.group_id
  option      = "unbind"
}

`
