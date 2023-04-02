package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsCompareTaskStopOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsCompareTaskStopOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_compare_task_stop_operation.compare_task_stop_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_compare_task_stop_operation.compare_task_stop_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsCompareTaskStopOperation = `

resource "tencentcloud_dts_compare_task_stop_operation" "compare_task_stop_operation" {
  job_id = "dts-8yv4w2i1"
  compare_task_id = "dts-8yv4w2i1-cmp-37skmii9"
}

`
