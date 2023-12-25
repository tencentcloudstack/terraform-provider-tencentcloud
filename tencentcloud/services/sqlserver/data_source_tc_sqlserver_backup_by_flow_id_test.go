package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverBackupByFlowIdDataSource_basic -v
func TestAccTencentCloudSqlserverBackupByFlowIdDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverGeneralBackupDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverBackupByFlowIdDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_backup_by_flow_id.example"),
				),
			},
		},
	})
}

const testAccSqlserverBackupByFlowIdDataSource = testAccSqlserverGeneralBackup + `
data "tencentcloud_sqlserver_backup_by_flow_id" "example" {
  instance_id = tencentcloud_sqlserver_general_backup.example.instance_id
  flow_id     = tencentcloud_sqlserver_general_backup.example.flow_id
}
`
