package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseDiskBackupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseDiskBackup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_disk_backup.disk_backup", "id"),
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_disk_backup.disk_backup", "disk_backup_name", "disk-backup"),
				),
			},
			{
				Config: testAccLighthouseDiskBackupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_disk_backup.disk_backup", "disk_backup_name", "disk-backup-update"),
				),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_disk_backup.disk_backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseDiskBackup = DefaultLighthoustVariables + `
resource "tencentcloud_lighthouse_disk_backup" "disk_backup" {
	disk_id = var.lighthouse_disk_id
	disk_backup_name = "disk-backup"
}
`

const testAccLighthouseDiskBackupUpdate = DefaultLighthoustVariables + `
resource "tencentcloud_lighthouse_disk_backup" "disk_backup" {
	disk_id = var.lighthouse_disk_id
	disk_backup_name = "disk-backup-update"
}
`
