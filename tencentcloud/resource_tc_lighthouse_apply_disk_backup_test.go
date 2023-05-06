package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseApplyDiskBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseApplyDiskBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_apply_disk_backup.apply_disk_backup", "id")),
			},
		},
	})
}

const testAccLighthouseApplyDiskBackup = DefaultLighthoustVariables + `

resource "tencentcloud_lighthouse_apply_disk_backup" "apply_disk_backup" {
  disk_id = var.lighthouse_backup_disk_id
  disk_backup_id = var.lighthouse_backup_id
}

`
