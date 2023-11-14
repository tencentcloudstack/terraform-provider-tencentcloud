package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbDatabasesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbDatabasesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_databases.databases")),
			},
		},
	})
}

const testAccCdbDatabasesDataSource = `

data "tencentcloud_cdb_databases" "databases" {
  instance_id = "cdb-c1nl9rpv"
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  database_regexp = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items = &lt;nil&gt;
  database_list {
		database_name = &lt;nil&gt;
		character_set = &lt;nil&gt;

  }
}

`
