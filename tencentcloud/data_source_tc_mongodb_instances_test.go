package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudMongodbInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstancesDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_instance.mongodb"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances", "instance_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances", "instance_list.0.instance_name", "tf-mongodb-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances", "instance_list.0.project_id", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances", "instance_list.0.cluster_type", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances", "instance_list.0.zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.vip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.vport"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.mongo_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.cpu"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.machine_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances", "instance_list.0.shard_quantity"),
				),
			},
		},
	})
}

const testAccMongodbInstancesDataSource = `
resource "tencentcloud_mongodb_instance" "mongodb_instance" {
	instance_name = "tf-mongodb-test"
	memory = 4
	volume = 100
	engine_version = "MONGO_36_WT"
	machine_type = "TGIO"
	available_zone = "ap-guangzhou-3"
	project_id = 0
	password = "test1234"
}

data "tencentcloud_mongodb_instances" "mongodb_instances" {
	instance_id = "${tencentcloud_mongodb_instance.mongodb_instance.id}"
}

`
