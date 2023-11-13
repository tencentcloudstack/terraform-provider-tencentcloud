package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbBackupTablesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbBackupTablesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_backup_tables.backup_tables")),
			},
		},
	})
}

const testAccCdbBackupTablesDataSource = `

data "tencentcloud_cdb_backup_tables" "backup_tables" {
  instance_id = "cdb-c1nl9rpv"
  start_time = "2022-07-12 10:29:20"
  database_name = &lt;nil&gt;
  search_table = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items {
		table_name = &lt;nil&gt;

  }
}

`
