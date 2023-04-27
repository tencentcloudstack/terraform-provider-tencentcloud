package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbDatabasesDataSource_basic -v
func TestAccTencentCloudMariadbDatabasesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbDatabasesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_databases.databases"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_databases.databases", "databases.#"),
				),
			},
		},
	})
}

const testAccMariadbDatabasesDataSource = testAccMariadbHourDbInstance + `

data "tencentcloud_mariadb_databases" "databases" {
  instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
}

`
