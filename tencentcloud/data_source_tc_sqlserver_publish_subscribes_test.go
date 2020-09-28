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
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.publish_instance_id", "mssql-82lhybgn"),
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.subscribe_instance_id", "mssql-12a60qdd"),
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.publish_subscribe_name", "example"),
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.publish_instance_ip", "10.1.0.17"),
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.subscribe_instance_ip", "10.1.0.11"),
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.publish_instance_name", "pub-keep"),
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.subscribe_instance_name", "sub-keep"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.publish_subscribe_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.publish_subscribes", "publish_subscribe_list.0.database_tuples.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.publish_foos", "publish_subscribe_list.0.publish_instance_id", "mssql-82lhybgn"),
					resource.TestCheckResourceAttr("data.tencentcloud_sqlserver_publish_subscribes.subscribe_foos", "publish_subscribe_list.0.subscribe_instance_id", "mssql-12a60qdd"),
				),
			},
		},
	})
}

const testAccTencentCloudSqlServerPublishSubscribeDataSourceConfig = `
resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = "mssql-82lhybgn"
	subscribe_instance_id           = "mssql-12a60qdd"
	publish_subscribe_name          = "example"
	database_tuples {
		publish_database            = "db_test_name"
		subscribe_database          = "db_test_name"
	}
}

data "tencentcloud_sqlserver_publish_subscribes" "publish_subscribes" {
	instance_id						= tencentcloud_sqlserver_publish_subscribe.example.publish_instance_id
	pub_or_sub_instance_id			= tencentcloud_sqlserver_publish_subscribe.example.subscribe_instance_id
	publish_subscribe_name          = tencentcloud_sqlserver_publish_subscribe.example.publish_subscribe_name
}

data "tencentcloud_sqlserver_publish_subscribes" "publish_foos" {
	instance_id				        = tencentcloud_sqlserver_publish_subscribe.example.publish_instance_id
}

data "tencentcloud_sqlserver_publish_subscribes" "subscribe_foos" {
	instance_id		                = tencentcloud_sqlserver_publish_subscribe.example.subscribe_instance_id
}
`
