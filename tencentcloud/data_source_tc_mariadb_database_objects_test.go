package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbDatabaseObjectsDataSource_basic -v
func TestAccTencentCloudMariadbDatabaseObjectsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbDatabaseObjectsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_database_objects.database_objects"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_database_objects.database_objects", "tables.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_database_objects.database_objects", "procs.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_database_objects.database_objects", "views.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_database_objects.database_objects", "funcs.#"),
				),
			},
		},
	})
}

const testAccMariadbDatabaseObjectsDataSource = testAccMariadbHourDbInstance + `

data "tencentcloud_mariadb_databases" "databases" {
  instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
}
  
data "tencentcloud_mariadb_database_objects" "database_objects" {
  instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
  db_name = data.tencentcloud_mariadb_databases.databases.databases[0].db_name
}

`
