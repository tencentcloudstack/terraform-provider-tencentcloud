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
resource "tencentcloud_dlc_attach_data_mask_policy" "example" {
  data_mask_strategy_policy_set {
    policy_info {
      database    = ""
      catalog     = ""
      table       = ""
      operation   = ""
      policy_type = ""
      column      = ""
      mode        = ""
    }

    data_mask_strategy_id = ""
    column_type           = ""
  }
}
`
