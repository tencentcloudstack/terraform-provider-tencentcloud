package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsMigrateCheckJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateCheckJob,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_check_job.migrate_check_job", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_migrate_check_job.migrate_check_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsMigrateCheckJob = `

resource "tencentcloud_dts_migrate_check_job" "migrate_check_job" {
  job_id = &lt;nil&gt;
        }

`
