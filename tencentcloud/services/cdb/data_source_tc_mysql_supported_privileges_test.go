package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlSupportedPrivilegesDataSource_basic -v
func TestAccTencentCloudMysqlSupportedPrivilegesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlSupportedPrivilegesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_supported_privileges.supported_privileges"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_supported_privileges.supported_privileges", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_supported_privileges.supported_privileges", "column_supported_privileges.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_supported_privileges.supported_privileges", "database_supported_privileges.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_supported_privileges.supported_privileges", "global_supported_privileges.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_supported_privileges.supported_privileges", "table_supported_privileges.#"),
				),
			},
		},
	})
}

const testAccMysqlSupportedPrivilegesDataSourceVar = `
variable "instance_id" {
	default = "` + tcacctest.DefaultDbBrainInstanceId + `"
}
`

const testAccMysqlSupportedPrivilegesDataSource = testAccMysqlSupportedPrivilegesDataSourceVar + `

data "tencentcloud_mysql_supported_privileges" "supported_privileges" {
  instance_id = var.instance_id
}

`
