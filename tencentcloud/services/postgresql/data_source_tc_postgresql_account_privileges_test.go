package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudPostgresqlAccountPrivilegesDataSource_basic -v
func TestAccTencentCloudPostgresqlAccountPrivilegesDataSource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlAccountPrivilegesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_account_privileges.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_account_privileges.example", "db_instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_account_privileges.example", "user_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_account_privileges.example", "database_object_set.#"),
				),
			},
		},
	})
}

const testAccPostgresqlAccountPrivilegesDataSource = `
data "tencentcloud_postgresql_account_privileges" "example" {
  db_instance_id = "postgres-3hk6b6tj"
  user_name      = "tf_example"
  database_object_set {
    object_name = "postgres"
    object_type = "database"
  }
}
`
