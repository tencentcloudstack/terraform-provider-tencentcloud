package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverUploadBackupInfoDataSource_basic -v
func TestAccTencentCloudSqlserverUploadBackupInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverUploadBackupInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_upload_backup_info.example")),
			},
		},
	})
}

const testAccSqlserverUploadBackupInfoDataSource = `
data "tencentcloud_sqlserver_upload_backup_info" "example" {
  instance_id         = "mssql-qelbzgwf"
  backup_migration_id = "mssql-backup-migration-8a0f3eht"
}
`
