package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlDatabasesDataSource_basic -v
func TestAccTencentCloudMysqlDatabasesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDatabasesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_databases.databases"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_databases.databases", "database_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_databases.databases", "database_list.0.character_set"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_databases.databases", "database_list.0.database_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_databases.databases", "items.#"),
				),
			},
		},
	})
}

const testAccMysqlDatabasesDataSourceVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultDbBrainInstanceId + `"
}
`

const testAccMysqlDatabasesDataSource = testAccMysqlDatabasesDataSourceVar + `

data "tencentcloud_mysql_databases" "databases" {
	instance_id = var.instance_id
	database_regexp = "^tf_ci_test"
}

`
