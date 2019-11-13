package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudMysqlInstanceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudMysqlInstanceDataSourceConfig(mysqlInstanceCommonTestCase),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.0.instance_name", "testAccMysql"),
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.0.pay_type", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.0.memory_size", "1000"),
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.0.volume_size", "50"),
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.0.engine_version", "5.7"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance.mysql", "instance_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance.mysql", "instance_list.0.subnet_id"),
				),
			},
		},
	})
}

func testAccTencentCloudMysqlInstanceDataSourceConfig(commonTestCase string) string {
	return fmt.Sprintf(`
%s
data "tencentcloud_mysql_instance" "mysql" {
	mysql_id = "${tencentcloud_mysql_instance.default.id}"
	instance_name = "${tencentcloud_mysql_instance.default.instance_name}"
}
	`, commonTestCase)
}
