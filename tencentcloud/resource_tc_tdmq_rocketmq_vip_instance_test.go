package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRocketmqVipInstanceResource_basic -v
func TestAccTencentCloudTdmqRocketmqVipInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
		},
		CheckDestroy: testAccCheckTdmqRocketmqVipInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqVipInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqVipInstanceExists("tencentcloud_tdmq_rocketmq_vip_instance.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_vip_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTdmqRocketmqVipInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqVipInstanceExists("tencentcloud_tdmq_rocketmq_vip_instance.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "id"),
				),
			},
		},
	})
}

func testAccCheckTdmqRocketmqVipInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TdmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rocketmq_vip_instance" {
			continue
		}

		res, err := service.DescribeTdmqRocketmqVipInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound.Instance" {
					return nil
				}
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tdmq_rocketmq_vip_instance %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTdmqRocketmqVipInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TdmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTdmqRocketmqVipInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tem application %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTdmqRocketmqVipInstance = `
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.1.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_tdmq_rocketmq_vip_instance" "example" {
  name         = "tx-example"
  spec         = "rocket-vip-basic-1"
  node_count   = 2
  storage_size = 200
  zone_ids     = [
    data.tencentcloud_availability_zones.zones.zones.0.id,
    data.tencentcloud_availability_zones.zones.zones.1.id
  ]

  vpc_info {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id

  }

  time_span = 1
}
`

const testAccTdmqRocketmqVipInstanceUpdate = `
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.1.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_tdmq_rocketmq_vip_instance" "example" {
  name         = "tx-example-update"
  spec         = "rocket-vip-basic-2"
  node_count   = 3
  storage_size = 600
  zone_ids     = [
    data.tencentcloud_availability_zones.zones.zones.0.id,
    data.tencentcloud_availability_zones.zones.zones.1.id
  ]

  vpc_info {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id

  }

  time_span = 1
}
`
