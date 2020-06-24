package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudSqlserverDB_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSqlserverDB,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_db.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_db.foo", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_db.foo", "name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_db.foo", "charset"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_db.foo", "remark"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_db.foo", "create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_db.foo", "status"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudSqlserverDB = testAccSqlserverDB_basic + `
data "tencentcloud_sqlserver_db" "foo" {
  instance_id = tencentcloud_sqlserver_db.mysqlserver_db.instance_id
  name        = tencentcloud_sqlserver_db.mysqlserver_db.name
}
`
