package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsSyncJobStopOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncJobStopOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_stop_operation.sync_job_stop_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_sync_job_stop_operation.sync_job_stop_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsSyncJobStopOperation = `

resource "tencentcloud_dts_sync_job_stop_operation" "sync_job_stop_operation" {
  job_id = "sync-werwfs23"
  new_instance_class = "large"
}

`
