package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsMigrateJobConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateJobConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.migrate_job_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_migrate_job_config.migrate_job_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsMigrateJobConfig = `

resource "tencentcloud_dts_migrate_job_config" "migrate_job_config" {
  job_id = "dts-ekmhr27i"
  complete_mode = "immediately"
  action = "dts-ekmhr27i"
}

`
