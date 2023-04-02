package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsSyncJobResumeOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncJobResumeOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_resume_operation.sync_job_resume_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_sync_job_resume_operation.sync_job_resume_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsSyncJobResumeOperation = `

resource "tencentcloud_dts_sync_job_resume_operation" "sync_job_resume_operation" {
  job_id = "sync-werwfs23"
}

`
