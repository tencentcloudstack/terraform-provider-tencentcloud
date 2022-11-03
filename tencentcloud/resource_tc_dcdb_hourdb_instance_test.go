package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudDcdbHourdbInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcdbHourdbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbHourdbInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcdbHourdbInstanceExists("tencentcloud_dcdb_hourdb_instance.hourdb_instance"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "zones.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "security_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "instance_name", "test_dcdc_dc_instance"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "shard_memory", "2"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "shard_storage", "10"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "shard_node_count", "2"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "shard_count", "2"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "db_version_id", "8.0"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "resource_tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "resource_tags.0.tag_key", "aaa"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "resource_tags.0.tag_value", "bbb"),
				),
			},
			{
				Config: testAccDcdbHourdbInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcdbHourdbInstanceExists("tencentcloud_dcdb_hourdb_instance.hourdb_instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "security_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "instance_name", "test_dcdc_dc_instance_CHANGED"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_hourdb_instance.hourdb_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDcdbHourdbInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	dcdbService := DcdbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dcdb_hourdb_instance" {
			continue
		}

		ret, err := dcdbService.DescribeDcdbHourdbInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if *ret.TotalCount > 0 || len(ret.Instances) > 0 {
			return fmt.Errorf("dcdb hourdb instance still exist, instanceId: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDcdbHourdbInstanceExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("dcdb hourdb instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("dcdb hourdb instance id is not set")
		}

		dcdbService := DcdbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		ret, err := dcdbService.DescribeDcdbHourdbInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if *ret.TotalCount == 0 || len(ret.Instances) == 0 {
			return fmt.Errorf("dcdb hourdb instance not found, instanceId: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccDcdbHourdb_vpc_config = `
data "tencentcloud_security_groups" "internal" {
	name = "default"
}

data "tencentcloud_vpc_instances" "vpc" {
	name ="Default-VPC"
}
	
data "tencentcloud_vpc_subnets" "subnet" {
	vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
}
	
locals {
	vpc_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.vpc_id
	subnet_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id
	sg_id = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
}
`

const testAccDcdbHourdbInstance_basic = testAccDcdbHourdb_vpc_config + `

resource "tencentcloud_dcdb_hourdb_instance" "hourdb_instance" {
  instance_name = "test_dcdc_dc_instance"
  zones = ["ap-guangzhou-5"]
  shard_memory = "2"
  shard_storage = "10"
  shard_node_count = "2"
  shard_count = "2"
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
  security_group_id = local.sg_id
  db_version_id = "8.0"
  resource_tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
}

`

const testAccDcdbHourdbInstance_update = testAccDcdbHourdb_vpc_config + `

resource "tencentcloud_dcdb_hourdb_instance" "hourdb_instance" {
  instance_name = "test_dcdc_dc_instance_CHANGED"
  zones = ["ap-guangzhou-5"]
  shard_memory = "2"
  shard_storage = "10"
  shard_node_count = "2"
  shard_count = "2"
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
  security_group_id = local.sg_id
  db_version_id = "8.0"
  resource_tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
}

`
