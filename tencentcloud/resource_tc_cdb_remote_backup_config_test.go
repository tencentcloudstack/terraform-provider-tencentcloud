package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbRemoteBackupConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbRemoteBackupConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_remote_backup_config.remote_backup_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_remote_backup_config.remote_backup_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbRemoteBackupConfig = `

resource "tencentcloud_cdb_remote_backup_config" "remote_backup_config" {
  instance_id = "cdb-c1nl9rpv"
  remote_backup_save = "on"
  remote_binlog_save = "on"
  remote_region = 
  expire_days = 7
}

`
