package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCbsDiskBackupRollbackOperationResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsDiskBackupRollbackOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_disk_backup_rollback_operation.operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cbs_disk_backup_rollback_operation.operation", "is_rollback_completed", "true"),
				),
			},
		},
	})
}

const testAccCbsDiskBackupRollbackOperation = CbsBackUp + `
resource "tencentcloud_cbs_disk_backup" "disk_backup" {
	disk_id = var.cbs_backup_disk_id
	disk_backup_name = "test-backup" 
}

resource "tencentcloud_cbs_disk_backup_rollback_operation" "operation" {
disk_backup_id  = tencentcloud_cbs_disk_backup.disk_backup.id
disk_id = var.cbs_backup_disk_id
}
`
