package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbDatabaseTableDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbDatabaseTableDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_database_table.database_table")),
			},
		},
	})
}

const testAccDcdbDatabaseTableDataSource = `

data "tencentcloud_dcdb_database_table" "database_table" {
  instance_id = "dcdbt-ow7t8lmc"
  db_name = &lt;nil&gt;
  table = &lt;nil&gt;
  }

`
