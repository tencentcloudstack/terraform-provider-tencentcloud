package mongodb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMongodbInstanceBackupRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceBackupRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_backup_rule.backup_rule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_backup_rule.backup_rule", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_backup_rule.backup_rule", "backup_method"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_backup_rule.backup_rule", "backup_time", "10"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_backup_rule.backup_rule", "backup_retention_period"),
				),
			},
			{
				Config: testAccMongodbInstanceBackupRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_backup_rule.backup_rule", "backup_time", "14"),
				),
			},
			{
				ResourceName: "tencentcloud_mongodb_instance_backup_rule.backup_rule",
				ImportState:  true,
			},
		},
	})
}

const testAccMongodbInstanceBackupRule = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-encryption-test"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_mongodb_instance_backup_rule" "backup_rule" {
    instance_id = tencentcloud_mongodb_instance.mongodb.id
    backup_method = 0
    backup_time = 10
}
`

const testAccMongodbInstanceBackupRuleUpdate = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-encryption-test"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_mongodb_instance_backup_rule" "backup_rule" {
    instance_id = tencentcloud_mongodb_instance.mongodb.id
    backup_method = 0
    backup_time = 14
}
`
