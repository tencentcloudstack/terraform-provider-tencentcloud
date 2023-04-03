package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDtsMigrateJobResumeOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateJobResumeOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_resume_operation.migrate_job_resume_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_migrate_job_resume_operation.migrate_job_resume_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsMigrateJobResumeOperation = `

resource "tencentcloud_dts_migrate_job_resume_operation" "migrate_job_resume_operation" {
  job_id = "dts-ekmhr27i"
  resume_option = "normal"
}

resource "tencentcloud_dts_migrate_job_resume_operation" "resume" {
	job_id = "%"
	resume_option = "normal"
  }

`
