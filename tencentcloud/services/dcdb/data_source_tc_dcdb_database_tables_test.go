package dcdb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDCDBDatabaseTablesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDCDBDatabaseTablesDataSource, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_database_tables.database_tables"),
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
