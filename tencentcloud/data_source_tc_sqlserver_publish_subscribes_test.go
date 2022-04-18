package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudSqlserverPublishSubscribeDataSource(t *testing.T) {
	t.Parallel()
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

const testAccTencentCloudSqlServerPublishSubscribeDataSourceConfig = defaultVpcSubnets + defaultSecurityGroupData + CommonPresetSQLServer + `

resource "tencentcloud_sqlserver_instance" "publish_instance" {
  name                  = "tf_sqlserver_publish_instance"
  availability_zone         = var.default_az
  charge_type               = "POSTPAID_BY_HOUR"
  vpc_id                    = local.vpc_id
  subnet_id                 = local.subnet_id
  project_id            = 0
  memory                = 2
  storage               = 10
  maintenance_week_set  = [1,2,3]
  maintenance_start_time= "09:00"
  maintenance_time_span = 3
  security_groups       = [local.sg_id]
}

resource "tencentcloud_sqlserver_instance" "subscribe_instance" {
  name                      = "tf_sqlserver_subscribe_instance"
  availability_zone         = var.default_az
  charge_type               = "POSTPAID_BY_HOUR"
  vpc_id                    = local.vpc_id
  subnet_id                 = local.subnet_id
  memory                    = 2
  storage                   = 10
  maintenance_week_set      = [1,2,3]
  maintenance_start_time    = "09:00"
  maintenance_time_span     = 3
  security_groups           = [local.sg_id]
}

resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = local.sqlserver_id
	subscribe_instance_id           = tencentcloud_sqlserver_instance.subscribe_instance.id
	publish_subscribe_name          = "example"
	delete_subscribe_db             = false
	database_tuples {
		publish_database            = local.sqlserver_db
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
