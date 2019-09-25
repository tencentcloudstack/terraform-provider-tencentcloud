package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudMongodbShardingInstance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbShardingInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_sharding_instance.mongodb"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "instance_name", "tf-mongodb-sharding"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "shard_quantity", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "nodes_per_shard", "3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "memory", "4"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "volume", "100"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "engine_version", "MONGO_36_WT"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "machine_type", "TGIO"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "available_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "tags.test", "test"),
				),
			},
			{
				Config: testAccMongodbShardingInstanceUpdateTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_sharding_instance.mongodb"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "tags.abc", "abc"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_sharding_instance.mongodb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_groups", "password"},
			},
		},
	})
}

func testAccCheckMongodbShardingInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mongodbService := MongodbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mongodb_sharding_instance" {
			continue
		}

		_, err := mongodbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("mongodb sharding instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccMongodbShardingInstance = `
resource "tencentcloud_mongodb_sharding_instance" "mongodb" {
  instance_name   = "tf-mongodb-sharding"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_36_WT"
  machine_type    = "TGIO"
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234"

  tags = {
    "test" = "test"
  }
}
`

const testAccMongodbShardingInstanceUpdateTags = `
resource "tencentcloud_mongodb_sharding_instance" "mongodb" {
  instance_name   = "tf-mongodb-sharding"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_36_WT"
  machine_type    = "TGIO"
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234"

  tags = {
    "abc" = "abc"
  }
}
`
