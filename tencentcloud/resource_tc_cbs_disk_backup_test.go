package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCbsDiskBackupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsDiskBackup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_disk_backup.disk_backup", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cbs_disk_backup.disk_backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCbsDiskBackup = `
resource "tencentcloud_cbs_disk_backup" "disk_backup" {
  disk_id = "disk-r69pg9vw"
  disk_backup_name = "test-disk-backup"
}
`
