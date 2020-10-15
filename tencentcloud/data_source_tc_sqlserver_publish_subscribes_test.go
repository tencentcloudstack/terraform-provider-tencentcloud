package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudSqlserverPublishSubscribeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverPublishSubscribeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudSqlServerPublishSubscribeDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverPublishSubscribeExists("tencentcloud_sqlserver_publish_subscribe.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.publish_instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.subscribe_instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.publish_subscribe_name", "example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.publish_instance_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.subscribe_instance_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.publish_instance_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.subscribe_instance_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.publish_subscribe_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.database_tuples.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_publish_subscribes.publish_foos", "publish_subscribe_list.0.publish_instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_publish_subscribes.subscribe_foos", "publish_subscribe_list.0.subscribe_instance_id"),
				),
			},
		},
	})
}

const testAccTencentCloudSqlServerPublishSubscribeDataSourceConfig = testAccSqlserverInstanceBasic + `
resource "tencentcloud_security_group" "foo" {
  name = "test-sg-tf"
}

resource "tencentcloud_sqlserver_instance" "publish_instance" {
  name                  = "tf_sqlserver_publish_instance"
  availability_zone     = var.availability_zone
  charge_type           = "POSTPAID_BY_HOUR"
  vpc_id                = "` + defaultVpcId + `"
  subnet_id             = "` + defaultSubnetId + `"
  project_id            = 0
  memory                = 2
  storage               = 10
  maintenance_week_set  = [1,2,3]
  maintenance_start_time= "09:00"
  maintenance_time_span = 3
  security_groups       = [tencentcloud_security_group.foo.name]
}

resource "tencentcloud_sqlserver_instance" "subscribe_instance" {
  name                      = "tf_sqlserver_subscribe_instance"
  availability_zone         = var.availability_zone
  charge_type               = "POSTPAID_BY_HOUR"
  vpc_id                    = "` + defaultVpcId + `"
  subnet_id                 = "` + defaultSubnetId + `"
  project_id                = 0
  memory                    = 2
  storage                   = 10
  maintenance_week_set      = [1,2,3]
  maintenance_start_time    = "09:00"
  maintenance_time_span     = 3
  security_groups           = [tencentcloud_security_group.foo.name]
}

resource "tencentcloud_sqlserver_db" "test_publish_subscribe" {
  instance_id   = tencentcloud_sqlserver_instance.publish_instance.id
  name          = "test111"
  charset       = "Chinese_PRC_BIN"
  remark        = "testACC-remark"
}

resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = tencentcloud_sqlserver_instance.publish_instance.id
	subscribe_instance_id           = tencentcloud_sqlserver_instance.subscribe_instance.id
	publish_subscribe_name          = "example"
	delete_subscribe_db             = false
	database_tuples {
		publish_database            = tencentcloud_sqlserver_db.test_publish_subscribe.name
	}
}

data "tencentcloud_sqlserver_publish_subscribes" "publish_subscribes" {
	instance_id                     = tencentcloud_sqlserver_publish_subscribe.example.publish_instance_id
	pub_or_sub_instance_id          = tencentcloud_sqlserver_publish_subscribe.example.subscribe_instance_id
	publish_subscribe_name          = tencentcloud_sqlserver_publish_subscribe.example.publish_subscribe_name
}

data "tencentcloud_sqlserver_publish_subscribes" "publish_foos" {
	instance_id				        = tencentcloud_sqlserver_publish_subscribe.example.publish_instance_id
}

data "tencentcloud_sqlserver_publish_subscribes" "subscribe_foos" {
	instance_id		                = tencentcloud_sqlserver_publish_subscribe.example.subscribe_instance_id
}
`
