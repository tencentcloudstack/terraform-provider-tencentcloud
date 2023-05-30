package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlRemoteBackupConfigResource_basic -v
func TestAccTencentCloudMysqlRemoteBackupConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRemoteBackupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_remote_backup_config.remote_backup_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_remote_backup_config.remote_backup_config", "remote_backup_save", "on"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_remote_backup_config.remote_backup_config", "remote_binlog_save", "on"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_remote_backup_config.remote_backup_config", "remote_region.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_remote_backup_config.remote_backup_config", "remote_region.0", "ap-shanghai"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_remote_backup_config.remote_backup_config", "expire_days", "7"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_remote_backup_config.remote_backup_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlRemoteBackupConfigVar = `
variable "instance_id" {
  default = "` + defaultDbBrainInstanceId + `"
}
`

const testAccMysqlRemoteBackupConfig = testAccMysqlRemoteBackupConfigVar + `

resource "tencentcloud_mysql_remote_backup_config" "remote_backup_config" {
	instance_id = var.instance_id
	remote_backup_save = "on"
	remote_binlog_save = "on"
	remote_region = ["ap-shanghai"]
	expire_days = 7
}

`
