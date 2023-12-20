package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlInstTablesDataSource_basic -v
func TestAccTencentCloudMysqlInstTablesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlInstTablesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_inst_tables.inst_tables"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_inst_tables.inst_tables", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_inst_tables.inst_tables", "items.#"),
				),
			},
		},
	})
}

const testAccMysqlInstTablesDataSourceVar = `
variable "instance_id" {
	default = "` + tcacctest.DefaultDbBrainInstanceId + `"
}
`

const testAccMysqlInstTablesDataSource = testAccMysqlInstTablesDataSourceVar + `

data "tencentcloud_mysql_inst_tables" "inst_tables" {
	instance_id = var.instance_id
	database = "tf_ci_test"
	# table_regexp = ""
}

`
