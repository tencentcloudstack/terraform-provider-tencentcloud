package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixPostgresqlDeleteLogBackupOperationResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlDeleteLogBackupOperation,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_delete_log_backup_operation.delete_log_backup_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_delete_log_backup_operation.delete_log_backup_operation", "db_instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_delete_log_backup_operation.delete_log_backup_operation", "log_backup_id"),
				),
			},
		},
	})
}

const testAccPostgresqlDeleteLogBackupOperation = tcacctest.CommonPresetPGSQL + `

resource "tencentcloud_postgresql_delete_log_backup_operation" "delete_log_backup_operation" {
  db_instance_id = local.pgsql_id
  log_backup_id = "9e93596c-c5b1-557e-aa87-8b857d79e283" 
  # use the data source after tencentcloud_postgresql_base_backups ready
  #log_backup_id = data.tencentcloud_postgresql_base_backups.foo.log_backup_set.0.id
}

`
