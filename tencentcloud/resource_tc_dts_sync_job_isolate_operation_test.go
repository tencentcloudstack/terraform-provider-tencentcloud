package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsSyncJobIsolateOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncJobIsolateOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_isolate_operation.sync_job_isolate_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_sync_job_isolate_operation.sync_job_isolate_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsSyncJobIsolateOperation = `

resource "tencentcloud_dts_sync_job_isolate_operation" "sync_job_isolate_operation" {
  job_id = "sync-werwfs23"
  new_instance_class = "large"
}

`
