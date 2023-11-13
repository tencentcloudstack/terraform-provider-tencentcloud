package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbTablesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbTablesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_tables.tables")),
			},
		},
	})
}

const testAccCdbTablesDataSource = `

data "tencentcloud_cdb_tables" "tables" {
  instance_id = "cdb-c1nl9rpv"
  database = &lt;nil&gt;
  offset = 0
  limit = 20
  table_regexp = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items = &lt;nil&gt;
}

`
