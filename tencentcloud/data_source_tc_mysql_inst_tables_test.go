package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlInstTablesDataSource_basic -v
func TestAccTencentCloudMysqlInstTablesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlInstTablesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_inst_tables.inst_tables"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_inst_tables.inst_tables", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_inst_tables.inst_tables", "items.#"),
				),
			},
		},
	})
}

const testAccMysqlInstTablesDataSourceVar = `
variable "instance_id" {
	default = "` + defaultDbBrainInstanceId + `"
}
`

const testAccMysqlInstTablesDataSource = testAccMysqlInstTablesDataSourceVar + `

data "tencentcloud_mysql_inst_tables" "inst_tables" {
	instance_id = var.instance_id
	database = "tf_ci_test"
	# table_regexp = ""
}

`
