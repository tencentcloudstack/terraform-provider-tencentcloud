package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlDeleteLogBackupOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlDeleteLogBackupOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgresql_delete_log_backup_operation.delete_log_backup_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgresql_delete_log_backup_operation.delete_log_backup_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresqlDeleteLogBackupOperation = `

resource "tencentcloud_postgresql_delete_log_backup_operation" "delete_log_backup_operation" {
  db_instance_id = ""
  log_backup_id = ""
}

`
