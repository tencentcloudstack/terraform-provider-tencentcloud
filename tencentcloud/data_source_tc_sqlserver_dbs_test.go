package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccDataSourceTencentCloudSqlserverDBs_basic -v
func TestAccDataSourceTencentCloudSqlserverDBs_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSqlserverDB,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_dbs.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_dbs.foo", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_dbs.foo", "db_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_dbs.foo", "db_list.0.charset"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_dbs.foo", "db_list.0.remark"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_dbs.foo", "db_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_dbs.foo", "db_list.0.status"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudSqlserverDB = CommonPresetSQLServer + `
data "tencentcloud_sqlserver_dbs" "foo" {
  instance_id = local.sqlserver_id
}
`
