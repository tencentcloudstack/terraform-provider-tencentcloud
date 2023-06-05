package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbDbSyncModeConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDcdbDbSyncModeConfig, defaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_sync_mode_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_sync_mode_config.config", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_sync_mode_config.config", "sync_mode", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDcdbDbSyncModeConfig_update, defaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_sync_mode_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_sync_mode_config.config", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_sync_mode_config.config", "sync_mode", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_db_sync_mode_config.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbDbSyncModeConfig = `

resource "tencentcloud_dcdb_db_sync_mode_config" "config" {
  instance_id = "%s"
  sync_mode = 2
}

`

const testAccDcdbDbSyncModeConfig_update = `

resource "tencentcloud_dcdb_db_sync_mode_config" "config" {
  instance_id = "%s"
  sync_mode = 1
}

`
