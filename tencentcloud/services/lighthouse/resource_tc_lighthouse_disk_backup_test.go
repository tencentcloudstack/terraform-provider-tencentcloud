package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseDiskBackupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
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

const testAccLighthouseDiskBackup = tcacctest.DefaultLighthoustVariables + `
resource "tencentcloud_lighthouse_disk_backup" "disk_backup" {
	disk_id = var.lighthouse_disk_id
	disk_backup_name = "disk-backup"
}
`

const testAccLighthouseDiskBackupUpdate = tcacctest.DefaultLighthoustVariables + `
resource "tencentcloud_lighthouse_disk_backup" "disk_backup" {
	disk_id = var.lighthouse_disk_id
	disk_backup_name = "disk-backup-update"
}
`
