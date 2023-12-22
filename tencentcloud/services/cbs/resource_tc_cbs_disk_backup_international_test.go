package cbs_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudInternationalCbsResource_diskBackup(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalCbsDiskBackup,
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

const testAccInternationalCbsDiskBackup = `
resource "tencentcloud_cbs_disk_backup" "disk_backup" {
	disk_id          = "disk-j8wrj3uq"
	disk_backup_name = "test-disk-backup"
}
`
