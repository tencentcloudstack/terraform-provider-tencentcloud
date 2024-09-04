package controlcenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixBatchApplyAccountBaselinesResource_basic -v
func TestAccTencentCloudNeedFixBatchApplyAccountBaselinesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchApplyAccountBaselines,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_batch_apply_account_baselines.example", "id"),
				),
			},
		},
	})
}

const testAccBatchApplyAccountBaselines = `
resource "tencentcloud_batch_apply_account_baselines" "example" {
  member_uin_list = [
    10037652245,
    10037652240,
  ]

  baseline_config_items {
    identifier    = "TCC-AF_SHARE_IMAGE"
    configuration = "{\"Images\":[{\"Region\":\"ap-guangzhou\",\"ImageId\":\"img-mcdsiqrx\",\"ImageName\":\"demo1\"}, {\"Region\":\"ap-guangzhou\",\"ImageId\":\"img-esxgkots\",\"ImageName\":\"demo2\"}]}"
  }
}
`
