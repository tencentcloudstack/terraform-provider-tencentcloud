package mongodb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmongodb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudMongodbStandbyInstanceResource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMongodbStandbyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbStandbyInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_standby_instance.mongodb"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "instance_name", "tf-mongodb-standby-test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "memory", "4"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "volume", "100"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb", "engine_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb", "machine_type"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "available_zone", "ap-guangzhou-4"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "charge_type", svcmongodb.MONGODB_CHARGE_TYPE_POSTPAID),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "prepaid_period"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb", "father_instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "father_instance_region", "ap-guangzhou"),
				),
			},
			{
				Config: testAccMongodbStandbyInstance_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "instance_name", "tf-mongodb-standby-test-update"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "memory", "8"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "volume", "200"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "tags.abc", "abc"),
				),
			},
			{
				Config: testAccMongodbStandbyInstance_securityGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "security_groups.0", "sg-05f7wnhn"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_standby_instance.mongodb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_groups", "auto_renew_flag", "password"},
			},
		},
	})
}

func testAccCheckMongodbStandbyInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mongodbService := svcmongodb.NewMongodbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mongodb_standby_instance" {
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

const testAccMongodbInstanceStandby = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "father_vpc" {
	name       = "mongodb-standby-father-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "father_subnet" {
	vpc_id            = tencentcloud_vpc.father_vpc.id
	name              = "mongodb-standby-father-subnet"
	cidr_block        = "10.0.1.0/24"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-test"
  memory         = 4
  volume         = 100
  engine_version  = local.sharding_engine_version
  machine_type    = local.sharding_machine_type
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = tencentcloud_vpc.father_vpc.id
  subnet_id      = tencentcloud_subnet.father_subnet.id

  tags = {
    test = "test"
  }
}
`

const testAccMongodbStandbyInstance = testAccMongodbInstanceStandby + `

resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.father_vpc.id
	name              = "mongodb-standby-subnet"
	cidr_block        = "10.0.2.0/24"
	availability_zone = "ap-guangzhou-4"
}

resource "tencentcloud_mongodb_standby_instance" "mongodb" {
  instance_name          = "tf-mongodb-standby-test"
  memory                 = 4
  volume                 = 100
  available_zone         = "ap-guangzhou-4"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.mongodb.id
  father_instance_region = "ap-guangzhou"
  vpc_id         = tencentcloud_vpc.father_vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
  security_groups  = [local.security_group_id]

  tags = {
    test = "test"
  }
}
`

const testAccMongodbStandbyInstance_update = testAccMongodbInstanceStandby + `
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.father_vpc.id
	name              = "mongodb-standby-subnet"
	cidr_block        = "10.0.2.0/24"
	availability_zone = "ap-guangzhou-4"
}

resource "tencentcloud_mongodb_standby_instance" "mongodb" {
  instance_name          = "tf-mongodb-standby-test-update"
  memory                 = 8
  volume                 = 200
  available_zone         = "ap-guangzhou-4"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.mongodb.id
  father_instance_region = "ap-guangzhou"
  vpc_id         = tencentcloud_vpc.father_vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
  security_groups  = [local.security_group_id]

  tags = {
    abc = "abc"
  }
}
`

const testAccMongodbStandbyInstance_securityGroup = testAccMongodbInstanceStandby + `
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.father_vpc.id
	name              = "mongodb-standby-subnet"
	cidr_block        = "10.0.2.0/24"
	availability_zone = "ap-guangzhou-4"
}

resource "tencentcloud_mongodb_standby_instance" "mongodb" {
  instance_name          = "tf-mongodb-standby-test-update"
  memory                 = 8
  volume                 = 200
  available_zone         = "ap-guangzhou-4"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.mongodb.id
  father_instance_region = "ap-guangzhou"
  vpc_id         = tencentcloud_vpc.father_vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
  security_groups = ["sg-05f7wnhn"]

  tags = {
    abc = "abc"
  }
}
`
