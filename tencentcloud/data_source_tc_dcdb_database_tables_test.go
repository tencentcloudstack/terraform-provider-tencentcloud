package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDCDBDatabaseTablesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDCDBDatabaseTablesDataSource, defaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_database_tables.database_tables"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_database_tables.database_tables", "cols.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_database_tables.database_tables", "db_name", "tf_test_db"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_database_tables.database_tables", "table", "tf_test_table"),
				),
			},
		},
	})
}

const testAccDCDBDatabaseTablesDataSource = `

data "tencentcloud_dcdb_database_tables" "database_tables" {
  instance_id = "%s"
  db_name = "tf_test_db"
  table = "tf_test_table"
}

`
