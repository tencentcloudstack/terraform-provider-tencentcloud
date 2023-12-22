package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverUploadIncrementalInfoDataSource_basic -v
func TestAccTencentCloudSqlserverUploadIncrementalInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverIncreBackupMigrationDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverUploadIncrementalInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_upload_incremental_info.example")),
			},
		},
	})
}

const testAccSqlserverUploadIncrementalInfoDataSource = `
data "tencentcloud_sqlserver_upload_incremental_info" "example" {
  instance_id              = "mssql-4tgeyeeh"
  backup_migration_id      = "mssql-backup-migration-83t5u3tv"
  incremental_migration_id = "mssql-incremental-migration-h36gkdxn"
}
`
