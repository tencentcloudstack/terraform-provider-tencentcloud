package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseBackupTablesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseBackupTablesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clickhouse_backup_tables.backup_tables")),
			},
		},
	})
}

const testAccClickhouseBackupTablesDataSource = DefaultClickhouseVariables + `
data "tencentcloud_clickhouse_backup_tables" "backup_tables" {
  instance_id = var.instance_id
  }
`
