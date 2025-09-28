package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcDataMaskStrategyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDataMaskStrategy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_data_mask_strategy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_data_mask_strategy.example", "strategy.#"),
				),
			},
			{
				Config: testAccDlcDataMaskStrategyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_data_mask_strategy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_data_mask_strategy.example", "strategy.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_dlc_data_mask_strategy.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcDataMaskStrategy = `
resource "tencentcloud_dlc_data_mask_strategy" "example" {
  strategy {
    strategy_name = "tf-example"
    strategy_desc = "description."
    groups {
      work_group_id = 70220
      strategy_type = "MASK"
    }
  }
}
`

const testAccDlcDataMaskStrategyUpdate = `
resource "tencentcloud_dlc_data_mask_strategy" "example" {
  strategy {
    strategy_name = "tf-example-update"
    strategy_desc = "description update."
    groups {
      work_group_id = 70219
      strategy_type = "MASK_NONE"
    }
  }
}
`
