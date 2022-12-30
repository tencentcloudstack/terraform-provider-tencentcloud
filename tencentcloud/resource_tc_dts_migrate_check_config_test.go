package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsMigrateCheckConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateCheckConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_check_config.migrate_check_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_migrate_check_config.migrate_check_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsMigrateCheckConfig = `

resource "tencentcloud_dts_migrate_check_config" "migrate_check_config" {
  job_id = ""
}

`
