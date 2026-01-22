package mongodb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMongodbReadOnlyInstanceResource_replset(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbReplsetReadOnlyInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_readonly_instance.mongodb"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_readonly_instance.mongodb", "instance_name", "tf-mongodb-readonly-test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_readonly_instance.mongodb", "memory", "4"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_readonly_instance.mongodb", "volume", "100"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_readonly_instance.mongodb", "cluster_type", "REPLSET"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_readonly_instance.mongodb",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"security_groups", "auto_renew_flag", "password"},
			},
		},
	})
}

func TestAccTencentCloudMongodbReadOnlyInstanceResource_shard(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbShardingReadOnlyInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_readonly_instance.sharding_mongodb"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_readonly_instance.sharding_mongodb", "instance_name", "tf-mongodb-readonly-shard"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_readonly_instance.sharding_mongodb", "memory", "4"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_readonly_instance.sharding_mongodb", "volume", "100"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_readonly_instance.sharding_mongodb", "cluster_type", "SHARD"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_readonly_instance.sharding_mongodb", "mongos_cpu", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_readonly_instance.sharding_mongodb", "mongos_memory", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_readonly_instance.sharding_mongodb", "mongos_node_num", "3"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_readonly_instance.sharding_mongodb",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"security_groups", "auto_renew_flag", "password"},
			},
		},
	})
}

const testAccMongodbReplsetReadOnlyInstance = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
  name       = "mongodb-sharding-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "mongodb-sharding-subnet"
  cidr_block        = "10.0.0.0/16"
  availability_zone = "ap-guangzhou-3"
}
resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-test"
  memory         = 4
  volume         = 100
  engine_version = local.engine_version
  machine_type   = local.machine_type
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_mongodb_readonly_instance" "mongodb" {
  instance_name          = "tf-mongodb-readonly-test"
  memory                 = 4
  volume                 = 100
  engine_version         = local.engine_version
  machine_type           = local.machine_type
  available_zone         = "ap-guangzhou-3"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.mongodb.id
  father_instance_region = "ap-guangzhou"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  security_groups        = [local.security_group_id]
  cluster_type           = "REPLSET"
}
`

const testAccMongodbShardingReadOnlyInstance = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-sharding-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-sharding-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_sharding_instance" "mongodb" {
  instance_name   = "tf-mongodb-sharding"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = local.sharding_memory
  volume          = local.sharding_volume
  engine_version  = local.sharding_engine_version
  machine_type    = local.sharding_machine_type
  security_groups = [local.security_group_id]
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234"
  mongos_cpu      = 1
  mongos_memory   = 2
  mongos_node_num = 3
  vpc_id    = tencentcloud_vpc.vpc.id
  subnet_id = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_mongodb_readonly_instance" "sharding_mongodb" {
  instance_name          = "tf-mongodb-readonly-shard"
  memory                 = 4
  volume                 = 100
  engine_version         = local.engine_version
  machine_type           = local.machine_type
  available_zone         = "ap-guangzhou-3"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_sharding_instance.mongodb.id
  father_instance_region = "ap-guangzhou"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  security_groups        = [local.security_group_id]
  cluster_type           = "SHARD"
  mongos_cpu             = 1
  mongos_memory          = 2
  mongos_node_num        = 3
}
`
