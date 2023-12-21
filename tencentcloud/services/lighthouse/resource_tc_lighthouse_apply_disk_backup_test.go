package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseApplyDiskBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseApplyDiskBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_apply_disk_backup.apply_disk_backup", "id")),
			},
		},
	})
}

const testAccLighthouseApplyDiskBackup = tcacctest.DefaultLighthoustVariables + `

resource "tencentcloud_lighthouse_apply_disk_backup" "apply_disk_backup" {
  disk_id = var.lighthouse_backup_disk_id
  disk_backup_id = var.lighthouse_backup_id
}

`
