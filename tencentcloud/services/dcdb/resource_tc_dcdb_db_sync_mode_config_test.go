package dcdb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbDbSyncModeConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		// PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDcdbDbSyncModeConfig, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_sync_mode_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_sync_mode_config.config", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_sync_mode_config.config", "sync_mode", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDcdbDbSyncModeConfig_update, tcacctest.DefaultDcdbInstanceId),
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
