package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbDatabasesDataSource_basic -v
func TestAccTencentCloudMariadbDatabasesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbDatabasesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_databases.databases"),
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
