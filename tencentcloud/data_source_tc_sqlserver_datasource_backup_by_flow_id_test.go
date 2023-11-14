package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverDatasourceBackupByFlowIdDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDatasourceBackupByFlowIdDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_datasource_backup_by_flow_id.datasource_backup_by_flow_id")),
			},
		},
	})
}

const testAccSqlserverDatasourceBackupByFlowIdDataSource = `

data "tencentcloud_sqlserver_datasource_backup_by_flow_id" "datasource_backup_by_flow_id" {
  instance_id = "mssql-i1z41iwd"
  flow_id = ""
                      }

`
