package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "memory", "8"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "volume", "200"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "engine_version", "MONGO_36_WT"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "machine_type", MONGODB_MACHINE_TYPE_HIO10G),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "available_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "charge_type", MONGODB_CHARGE_TYPE_POSTPAID),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "prepaid_period"),
				),
			},
			{
				Config: testAccMongodbShardingInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_sharding_instance.mongodb"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "instance_name", "tf-mongodb-sharding-update"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "memory", "4"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "volume", "100"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "tags.abc", "abc"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_sharding_instance.mongodb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_groups", "password", "auto_renew_flag"},
			},
			{
				Config: testAccMongodbShardingInstancePrepaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_sharding_instance.mongodb_prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "instance_name", "tf-mongodb-sharding-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "shard_quantity", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "nodes_per_shard", "3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "memory", "4"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "volume", "100"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "engine_version", "MONGO_40_WT"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "machine_type", MONGODB_MACHINE_TYPE_HIO10G),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "available_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "tags.test", "test-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "charge_type", MONGODB_CHARGE_TYPE_PREPAID),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "prepaid_period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "auto_renew_flag", "0"),
				),
			},
			{
				Config: testAccMongodbShardingInstancePrepaid_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_sharding_instance.mongodb_prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "instance_name", "tf-mongodb-sharding-prepaid-update"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "memory", "8"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "volume", "200"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "tags.prepaid", "prepaid"),
				),
			},
		},
	})
}

func testAccCheckMongodbShardingInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mongodbService := MongodbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mongodb_sharding_instance" {
			continue
		}

		_, has, err := mongodbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("mongodb instance still exists: %s", rs.Primary.ID)
	}
	return nil
}

const testAccMongodbShardingInstance = `
resource "tencentcloud_mongodb_sharding_instance" "mongodb" {
  instance_name   = "tf-mongodb-sharding"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 8
  volume          = 200
  engine_version  = "MONGO_36_WT"
  machine_type    = "TGIO"
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234"

  tags = {
    test = "test"
  }
}
`

const testAccMongodbShardingInstanceUpdate = `
resource "tencentcloud_mongodb_sharding_instance" "mongodb" {
  instance_name   = "tf-mongodb-sharding-update"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_36_WT"
  machine_type    = "TGIO"
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234update"

  tags = {
    abc = "abc"
  }
}
`

const testAccMongodbShardingInstancePrepaid = `
resource "tencentcloud_mongodb_sharding_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-sharding-prepaid"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_40_WT"
  machine_type    = "HIO10G"
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 0

  tags = {
    test = "test-prepaid"
  }
}
`

const testAccMongodbShardingInstancePrepaid_update = `
resource "tencentcloud_mongodb_sharding_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-sharding-prepaid-update"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 8
  volume          = 200
  engine_version  = "MONGO_40_WT"
  machine_type    = "HIO10G"
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234update"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 0

  tags = {
    prepaid = "prepaid"
  }
}
`
