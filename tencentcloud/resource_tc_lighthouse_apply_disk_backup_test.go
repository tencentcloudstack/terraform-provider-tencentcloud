package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseApplyDiskBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseApplyDiskBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_apply_disk_backup.apply_disk_backup", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_apply_disk_backup.apply_disk_backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseApplyDiskBackup = `

resource "tencentcloud_lighthouse_apply_disk_backup" "apply_disk_backup" {
  disk_id = "lhdisk-123456"
  disk_backup_id = "lhbak-1234556"
}

`
