package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbDatabaseTableDataSource_basic -v
func TestAccTencentCloudMariadbDatabaseTableDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbDatabaseTableDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_database_table.database_table"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_database_table.database_table", "cols.#"),
				),
			},
		},
	})
}

const testAccMariadbDatabaseTableDataSource = testAccMariadbHourDbInstance + `

data "tencentcloud_mariadb_databases" "databases" {
  instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
}
  
data "tencentcloud_mariadb_database_objects" "database_objects" {
  instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
  db_name = data.tencentcloud_mariadb_databases.databases.databases[0].db_name
}
  
data "tencentcloud_mariadb_database_table" "database_table" {
  instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
  db_name = data.tencentcloud_mariadb_databases.databases.databases[0].db_name
  table = data.tencentcloud_mariadb_database_objects.database_objects.tables[0].table
}

`
