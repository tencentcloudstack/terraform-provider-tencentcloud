package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudMongodbStandbyInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbStandbyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbStandbyInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_standby_instance.mongodb"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "instance_name", "tf-mongodb-standby-test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "memory", "4"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "volume", "100"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "engine_version", "MONGO_36_WT"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "machine_type", MONGODB_MACHINE_TYPE_HIO10G),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "available_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb", "charge_type", MONGODB_CHARGE_TYPE_POSTPAID),
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
				ResourceName:            "tencentcloud_mongodb_standby_instance.mongodb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_groups", "auto_renew_flag"},
			},
			{
				Config: testAccMongodbStandbyInstancePrepaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_standby_instance.mongodb_prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "instance_name", "tf-mongodb-standby-test-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "memory", "8"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "volume", "200"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "engine_version", "MONGO_40_WT"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "machine_type", MONGODB_MACHINE_TYPE_HIO10G),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "available_zone", "ap-guangzhou-4"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "tags.test", "test-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "charge_type", MONGODB_CHARGE_TYPE_PREPAID),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "prepaid_period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "auto_renew_flag", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "father_instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "father_instance_region", "ap-guangzhou"),
				),
			},
			{
				Config: testAccMongodbStandbyInstancePrepaid_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "instance_name", "tf-mongodb-standby-prepaid-update"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "memory", "4"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "volume", "100"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_standby_instance.mongodb_prepaid", "tags.prepaid", "prepaid"),
				),
			},
		},
	})
}

func testAccCheckMongodbStandbyInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mongodbService := MongodbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
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

const testAccMongodbInstanceHIO10G = `
resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-test"
  memory         = 4
  volume         = 100
  engine_version = "MONGO_36_WT"
  machine_type   = "HIO10G"
  available_zone = "ap-guangzhou-2"
  project_id     = 0
  password       = "test1234"

  tags = {
    test = "test"
  }
}
`

const testAccMongodbStandbyInstance = testAccMongodbInstanceHIO10G + `
resource "tencentcloud_mongodb_standby_instance" "mongodb" {
  instance_name          = "tf-mongodb-standby-test"
  memory                 = 4
  volume                 = 100
  available_zone         = "ap-guangzhou-3"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.mongodb.id
  father_instance_region = "ap-guangzhou"

  tags = {
    test = "test"
  }
}
`

const testAccMongodbStandbyInstance_update = testAccMongodbInstanceHIO10G + `
resource "tencentcloud_mongodb_standby_instance" "mongodb" {
  instance_name          = "tf-mongodb-standby-test-update"
  memory                 = 8
  volume                 = 200
  available_zone         = "ap-guangzhou-3"
  project_id             = 0
  father_instance_id     = tencentcloud_mongodb_instance.mongodb.id
  father_instance_region = "ap-guangzhou"

  tags = {
    abc = "abc"
  }
}
`

const testAccMongodbStandbyInstancePrepaid = testAccMongodbInstancePrepaid + `
resource "tencentcloud_mongodb_standby_instance" "mongodb_prepaid" {
  instance_name          = "tf-mongodb-standby-test-prepaid"
  memory                 = 8
  volume                 = 200
  available_zone         = "ap-guangzhou-4"
  project_id             = 0
  charge_type            = "PREPAID"
  prepaid_period         = 1
  auto_renew_flag        = 1
  father_instance_id     = tencentcloud_mongodb_instance.mongodb_prepaid.id
  father_instance_region = "ap-guangzhou"

  tags = {
    test = "test-prepaid"
  }
}
`

const testAccMongodbStandbyInstancePrepaid_update = testAccMongodbInstancePrepaid + `
resource "tencentcloud_mongodb_standby_instance" "mongodb_prepaid" {
  instance_name          = "tf-mongodb-standby-prepaid-update"
  memory                 = 4
  volume                 = 100
  available_zone         = "ap-guangzhou-4"
  project_id             = 0
  charge_type            = "PREPAID"
  prepaid_period         = 1
  auto_renew_flag        = 1
  father_instance_id     = tencentcloud_mongodb_instance.mongodb_prepaid.id
  father_instance_region = "ap-guangzhou"

  tags = {
    prepaid = "prepaid"
  }
}
`
