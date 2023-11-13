package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsMigrateJobOperateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateJobOperate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_operate.migrate_job_operate", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_migrate_job_operate.migrate_job_operate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsMigrateJobOperate = `

resource "tencentcloud_dts_migrate_job_operate" "migrate_job_operate" {
  job_id = "dts-ekmhr27i"
  complete_mode = "immediately"
}

`
