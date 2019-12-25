package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_instance.mongodb_instance"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.instance_name", "tf-mongodb-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.project_id", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.cluster_type", MONGODB_CLUSTER_TYPE_REPLSET),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.available_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.vpc_id", ""),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.subnet_id", ""),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.status", "2"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.vip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.vport"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.engine_version", "MONGO_36_WT"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.cpu"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.memory", "4"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.volume", "100"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.machine_type", "TGIO"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.shard_quantity"),
					resource.TestCheckResourceAttr("data.tencentcloud_mongodb_instances.mongodb_instances", "instance_list.0.tags.test", "test"),
				),
			},
		},
	})
}

const testAccMongodbInstancesDataSource = `
resource "tencentcloud_mongodb_instance" "mongodb_instance" {
  instance_name  = "tf-mongodb-test"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_36_WT"
  machine_type   = "TGIO"
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_mongodb_instances" "mongodb_instances" {
  instance_id = tencentcloud_mongodb_instance.mongodb_instance.id

  tags = tencentcloud_mongodb_instance.mongodb_instance.tags
}
`
