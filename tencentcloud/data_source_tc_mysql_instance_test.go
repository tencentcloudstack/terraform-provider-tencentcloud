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
				Config: testAccTencentCloudMysqlInstanceDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_mysql_instance.mysql", "instance_list.0.instance_name", "testAccTencentCloudMysqlInstanceDataSourceConfig"),
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

func testAccTencentCloudMysqlInstanceDataSourceConfig() string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_instance" "mysql" {
	pay_type = 1
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccTencentCloudMysqlInstanceDataSourceConfig"
	vpc_id = "vpc-fzdzrsir"
	subnet_id = "subnet-he8ldxx6"
	engine_version = "5.7"
	root_password = "test1234"
	availability_zone = "ap-guangzhou-4"
}

data "tencentcloud_mysql_instance" "mysql" {
	mysql_id = "${tencentcloud_mysql_instance.mysql.id}"
	instance_name = "${tencentcloud_mysql_instance.mysql.instance_name}"
}
	`)
}
