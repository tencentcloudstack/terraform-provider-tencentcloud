package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudPostgresqlRestoreDbInstanceObjectsOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlRestoreDbInstanceObjectsOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_restore_db_instance_objects_operation.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_restore_db_instance_objects_operation.example", "db_instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_restore_db_instance_objects_operation.example", "task_id"),
				),
			},
		},
	})
}

const testAccPostgresqlRestoreDbInstanceObjectsOperation = `
resource "tencentcloud_postgresql_restore_db_instance_objects_operation" "example" {
  db_instance_id  = "postgres-6bwgamo3"
  restore_objects = ["user"]
  backup_set_id   = "your-backup-set-id"
}
`
