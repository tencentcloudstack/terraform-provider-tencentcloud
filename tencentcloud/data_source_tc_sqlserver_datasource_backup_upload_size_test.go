package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverDatasourceBackupUploadSizeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDatasourceBackupUploadSizeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_datasource_backup_upload_size.datasource_backup_upload_size")),
			},
		},
	})
}

const testAccSqlserverDatasourceBackupUploadSizeDataSource = `

data "tencentcloud_sqlserver_datasource_backup_upload_size" "datasource_backup_upload_size" {
  instance_id = "mssql-i1z41iwd"
  backup_migration_id = ""
  incremental_migration_id = ""
  }

`
