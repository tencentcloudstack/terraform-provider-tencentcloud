package mongodb_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svcmongodb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"
)

func TestAccTencentCloudMongodbShardingInstanceResource_postpaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMongodbShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config:    testAccMongodbShardingInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_sharding_instance.mongodb"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "instance_name", "tf-mongodb-sharding"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "shard_quantity", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "nodes_per_shard", "3"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "memory"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "volume"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "engine_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "machine_type"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "available_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "charge_type", svcmongodb.MONGODB_CHARGE_TYPE_POSTPAID),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "prepaid_period"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "security_groups.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb", "hidden_zone"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "availability_zone_list.#", "3"),
				),
			},
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config:    testAccMongodbShardingInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_sharding_instance.mongodb"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "instance_name", "tf-mongodb-sharding-update"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb", "tags.abc", "abc"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_sharding_instance.mongodb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccTencentCloudMongodbShardingInstanceResource_prepaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMongodbShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
				Config:    testAccMongodbShardingInstancePrepaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_sharding_instance.mongodb_prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "instance_name", "tf-mongodb-sharding-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "shard_quantity", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "nodes_per_shard", "3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "available_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "tags.test", "test-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "charge_type", svcmongodb.MONGODB_CHARGE_TYPE_PREPAID),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "prepaid_period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "auto_renew_flag", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "hidden_zone"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "availability_zone_list.#", "3"),
				),
			},
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
				Config:    testAccMongodbShardingInstancePrepaid_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_sharding_instance.mongodb_prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "instance_name", "tf-mongodb-sharding-prepaid-update"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_sharding_instance.mongodb_prepaid", "tags.prepaid", "prepaid"),
				),
			},
		},
	})
}

func testAccCheckMongodbShardingInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mongodbService := svcmongodb.NewMongodbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccMongodbShardingInstance = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_mongodb_sharding_instance" "mongodb" {
  instance_name   = "tf-mongodb-sharding"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = local.sharding_memory
  volume          = local.sharding_volume
  engine_version  = local.sharding_engine_version
  machine_type    = local.sharding_machine_type
  security_groups  = [local.security_group_id]
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234"
  mongos_cpu = 1
  mongos_memory =  2
  mongos_node_num = 3
  tags = {
    test = "test"
  }
  security_groups = [var.sg_id]
  vpc_id = var.vpc_id
  subnet_id = var.subnet_id
  availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"]
  hidden_zone = "ap-guangzhou-6"
}
`

const testAccMongodbShardingInstanceUpdate = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_mongodb_sharding_instance" "mongodb" {
  instance_name   = "tf-mongodb-sharding-update"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = local.sharding_memory
  volume          = local.sharding_volume
  engine_version  = local.sharding_engine_version
  machine_type    = local.sharding_machine_type
  security_groups  = [local.security_group_id]
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234update"
  mongos_cpu = 1
  mongos_memory =  2
  mongos_node_num = 3

  tags = {
    abc = "abc"
  }
  security_groups = [var.sg_id]
  vpc_id = var.vpc_id
  subnet_id = var.subnet_id
  availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"]
  hidden_zone = "ap-guangzhou-6"
}
`

const testAccMongodbShardingInstancePrepaid = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_mongodb_sharding_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-sharding-prepaid"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = local.sharding_memory
  volume          = local.sharding_volume
  engine_version  = local.sharding_engine_version
  machine_type    = local.sharding_machine_type
  security_groups  = [local.security_group_id]
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 0
  mongos_cpu = 1
  mongos_memory =  2
  mongos_node_num = 3

  tags = {
    test = "test-prepaid"
  }
  vpc_id = var.vpc_id
  subnet_id = var.subnet_id
  availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"]
  hidden_zone = "ap-guangzhou-6"
}
`

const testAccMongodbShardingInstancePrepaid_update = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_mongodb_sharding_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-sharding-prepaid-update"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = local.sharding_memory
  volume          = local.sharding_volume
  engine_version  = local.sharding_engine_version
  machine_type    = local.sharding_machine_type
  security_groups  = [local.security_group_id]
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234update"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 0
  mongos_cpu = 1
  mongos_memory =  2
  mongos_node_num = 3

  tags = {
    prepaid = "prepaid"
  }
  vpc_id = var.vpc_id
  subnet_id = var.subnet_id
  availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"]
  hidden_zone = "ap-guangzhou-6"
}
`
