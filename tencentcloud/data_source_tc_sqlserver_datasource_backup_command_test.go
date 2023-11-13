package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverDatasourceBackupCommandDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDatasourceBackupCommandDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_datasource_backup_command.datasource_backup_command")),
			},
		},
	})
}

const testAccSqlserverDatasourceBackupCommandDataSource = `

data "tencentcloud_sqlserver_datasource_backup_command" "datasource_backup_command" {
  backup_file_type = "FULL"
  data_base_name = "db_name"
  is_recovery = "No"
  local_path = ""
  }

`
