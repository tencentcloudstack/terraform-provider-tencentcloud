package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsMigrateJobResumeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateJobResume,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_resume.migrate_job_resume", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_migrate_job_resume.migrate_job_resume",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsMigrateJobResume = `

resource "tencentcloud_dts_migrate_job_resume" "migrate_job_resume" {
  job_id = "dts-ekmhr27i"
  resume_option = "normal"
}

`
