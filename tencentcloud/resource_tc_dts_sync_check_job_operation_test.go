package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsSyncCheckJobOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncCheckJobOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_check_job_operation.sync_check_job_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_sync_check_job_operation.sync_check_job_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsSyncCheckJobOperation = `

resource "tencentcloud_dts_sync_check_job_operation" "sync_check_job_operation" {
  job_id = ""
  }

`
