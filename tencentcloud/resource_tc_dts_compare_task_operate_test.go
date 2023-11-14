package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsCompareTaskOperateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsCompareTaskOperate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_compare_task_operate.compare_task_operate", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_compare_task_operate.compare_task_operate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsCompareTaskOperate = `

resource "tencentcloud_dts_compare_task_operate" "compare_task_operate" {
  job_id = "dts-8yv4w2i1"
  compare_task_id = "dts-8yv4w2i1-cmp-37skmii9"
}

`
