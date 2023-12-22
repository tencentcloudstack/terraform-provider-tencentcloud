package cdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMysqlResource_instance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMysqlMasterInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalMysqlMasterInstance_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "instance_name", TestAccTencentCloudMysqlInstanceName),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_deploy_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_sync_mode", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "availability_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "first_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "second_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
				),
			},
			{
				Config: testAccInternationalMysqlMasterInstanceUp_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "instance_name", TestAccTencentCloudMysqlInstanceName),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_deploy_mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_sync_mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "availability_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "first_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone_4),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "second_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone_4),
				),
			},
		},
	})
}

func testAccInternationalMysqlMasterInstance_basic() string {
	return `
resource "tencentcloud_mysql_instance" "mysql_master" {
  charge_type       = "POSTPAID"
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  slave_deploy_mode = 0
  first_slave_zone  = "ap-guangzhou-3"
  second_slave_zone = "ap-guangzhou-3"
  slave_sync_mode   = 2
  force_delete      = true
}`
}

func testAccInternationalMysqlMasterInstanceUp_basic() string {
	return `
resource "tencentcloud_mysql_instance" "mysql_master" {
  charge_type       = "POSTPAID"
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  slave_deploy_mode = 1
  first_slave_zone  = "ap-guangzhou-4"
  second_slave_zone = "ap-guangzhou-4"
  slave_sync_mode   = 1
  force_delete      = true
}`
}
