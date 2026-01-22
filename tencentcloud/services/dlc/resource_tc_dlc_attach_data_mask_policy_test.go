package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcAttachDataMaskPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcAttachDataMaskPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_attach_data_mask_policy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_attach_data_mask_policy.example", "data_mask_strategy_policy_set.#"),
				),
			},
		},
	})
}

const testAccDlcAttachDataMaskPolicy = `
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

resource "tencentcloud_dlc_attach_data_mask_policy" "example" {
  data_mask_strategy_policy_set {
    policy_info {
      database    = "test"
      catalog     = "DataLakeCatalog"
      table       = "test"
      column      = "id"
    }

    data_mask_strategy_id = tencentcloud_dlc_data_mask_strategy.example.id
    column_type           = "string"
  }
}
`
