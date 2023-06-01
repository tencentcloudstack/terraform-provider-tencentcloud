package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsSyncJobContinueOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncJobContinueOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_continue_operation.sync_job_continue_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_sync_job_continue_operation.sync_job_continue_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsSyncJobContinueOperation = `

resource "tencentcloud_dts_sync_job_continue_operation" "sync_job_continue_operation" {
  job_id = "sync-werwfs23"
  new_instance_class = "large"
}

`
