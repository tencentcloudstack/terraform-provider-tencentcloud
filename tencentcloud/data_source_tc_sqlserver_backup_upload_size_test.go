package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverBackupUploadSizeDataSource_basic -v
func TestAccTencentCloudSqlserverBackupUploadSizeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverBackupUploadSizeDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_backup_upload_size.backup_upload_size"),
				),
			},
		},
	})
}

const testAccSqlserverBackupUploadSizeDataSource = `
data "tencentcloud_sqlserver_backup_upload_size" "example" {
  instance_id         = "mssql-4gmc5805"
  backup_migration_id = "mssql-backup-migration-9tj0sxnz"
}
`
