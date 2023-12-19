package cdb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlInstanceDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudMysqlInstanceDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.0.instance_name", tcacctest.DefaultMySQLName),
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.0.pay_type", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.0.memory_size", "4000"),
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.0.volume_size", "200"),
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.0.engine_version", "5.7"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance.mysql", "instance_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance.mysql", "instance_list.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance.mysql", "instance_list.0.auto_renew_flag"),
				),
			},
		},
	})
}

func testAccTencentCloudMysqlInstanceDataSourceConfig() string {
	return fmt.Sprintf(`
data "tencentcloud_mysql_instance" "mysql" {
	instance_name = "%s"
}
	`, tcacctest.DefaultMySQLName)
}
