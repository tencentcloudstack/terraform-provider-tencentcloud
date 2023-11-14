package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverUploadBackupInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverUploadBackupInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_upload_backup_info.upload_backup_info")),
			},
		},
	})
}

const testAccSqlserverUploadBackupInfoDataSource = `

data "tencentcloud_sqlserver_upload_backup_info" "upload_backup_info" {
  instance_id = "mssql-j8kv137v"
  backup_migration_id = "migration_id"
                }

`
