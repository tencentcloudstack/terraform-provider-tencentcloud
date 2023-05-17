package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverBackupByFlowIdDataSource_basic -v
func TestAccTencentCloudSqlserverBackupByFlowIdDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverGeneralBackupDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverBackupByFlowIdDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_backup_by_flow_id.backup_by_flow_id"),
				),
			},
		},
	})
}

const testAccSqlserverBackupByFlowIdDataSource = testAccSqlserverGeneralBackup + `
data "tencentcloud_sqlserver_backup_by_flow_id" "backup_by_flow_id" {
  instance_id = tencentcloud_sqlserver_general_backup.general_backup.instance_id
  flow_id = tencentcloud_sqlserver_general_backup.general_backup.flow_id
}
`
