package trocket_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRocketmqVipInstanceResource_basic -v -timeout=0
func TestAccTencentCloudTdmqRocketmqVipInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckTdmqRocketmqVipInstanceDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqVipInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqVipInstanceExists("tencentcloud_tdmq_rocketmq_vip_instance.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "spec"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "node_count"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "storage_size"),
					tcacctest.AccStepTimeSleepDuration(1*time.Minute),
				),
			},
			{
				Config: testAccTdmqRocketmqVipInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRocketmqVipInstanceExists("tencentcloud_tdmq_rocketmq_vip_instance.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "spec"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "node_count"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.example", "storage_size"),
					tcacctest.AccStepTimeSleepDuration(1*time.Minute),
				),
			},
		},
	})
}

func testAccCheckTdmqRocketmqVipInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctdmq.NewTdmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctdmq.NewTdmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
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
  ip_rules {
    ip_rule = "0.0.0.0/0"
    allow   = true
    remark  = "remark."
  }
}
`

const testAccTdmqRocketmqVipInstanceUpdate = `
data "tencentcloud_availability_zones" "zones" {}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_tdmq_rocketmq_vip_instance" "example" {
  name         = "tx-example-update"
  spec         = "rocket-vip-basic-2"
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
  ip_rules {
    ip_rule = "0.0.0.0/0"
    allow   = true
    remark  = "remark."
  }
}
`
