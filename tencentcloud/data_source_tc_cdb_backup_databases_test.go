package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbBackupDatabasesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbBackupDatabasesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_backup_databases.backup_databases")),
			},
		},
	})
}

const testAccCdbBackupDatabasesDataSource = `

data "tencentcloud_cdb_backup_databases" "backup_databases" {
  instance_id = "cdb-c1nl9rpv"
  start_time = "2022-07-12 10:29:20"
  search_database = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items {
		database_name = &lt;nil&gt;

  }
}

`
