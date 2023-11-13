package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCbsDiskBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsDiskBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cbs_disk_backup.disk_backup", "id")),
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
  disk_id = "disk-xxx"
  disk_backup_name = "xxx"
}

`
