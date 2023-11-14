package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverUploadIncrementalInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverUploadIncrementalInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_upload_incremental_info.upload_incremental_info")),
			},
		},
	})
}

const testAccSqlserverUploadIncrementalInfoDataSource = `

data "tencentcloud_sqlserver_upload_incremental_info" "upload_incremental_info" {
  instance_id = "mssql-j8kv137v"
  backup_migration_id = "migration_id"
  incremental_migration_id = ""
                }

`
