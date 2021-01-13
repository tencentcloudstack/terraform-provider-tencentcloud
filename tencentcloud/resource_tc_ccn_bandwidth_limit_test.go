package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const keyNameLimit1 = "tencentcloud_ccn_bandwidth_limit.limit1"

func TestAccTencentCloudCcnV3BandwidthLimitOuter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCcnBandwidthLimitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnOuterBandwidthLimitConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnBandwidthLimitExists(keyNameLimit1),
					resource.TestCheckResourceAttrSet(keyNameLimit1, "ccn_id"),
					resource.TestCheckResourceAttrSet(keyNameLimit1, "region"),
					resource.TestCheckResourceAttr(keyNameLimit1, "bandwidth_limit", "500"),
				),
			},
			{
				Config: testAccCcnOuterBandwidthLimitConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnBandwidthLimitExists(keyNameLimit1),
					resource.TestCheckResourceAttrSet(keyNameLimit1, "ccn_id"),
					resource.TestCheckResourceAttrSet(keyNameLimit1, "region"),
					resource.TestCheckResourceAttr(keyNameLimit1, "bandwidth_limit", "100"),
				),
			},
		},
	})
}

func TestAccTencentCloudCcnV3BandwidthLimitInter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCcnBandwidthLimitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnInterBandwidthLimitConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnBandwidthLimitExists(keyNameLimit1),
					resource.TestCheckResourceAttrSet(keyNameLimit1, "ccn_id"),
					resource.TestCheckResourceAttrSet(keyNameLimit1, "region"),
					resource.TestCheckResourceAttr(keyNameLimit1, "dst_region", "ap-beijing"),
					resource.TestCheckResourceAttr(keyNameLimit1, "bandwidth_limit", "500"),
				),
			},
			{
				Config: testAccCcnInterBandwidthLimitConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnBandwidthLimitExists(keyNameLimit1),
					resource.TestCheckResourceAttrSet(keyNameLimit1, "ccn_id"),
					resource.TestCheckResourceAttrSet(keyNameLimit1, "region"),
					resource.TestCheckResourceAttr(keyNameLimit1, "dst_region", "ap-nanjing"),
					resource.TestCheckResourceAttr(keyNameLimit1, "bandwidth_limit", "100"),
				),
			},
		},
	})
}

func testAccCheckCcnBandwidthLimitDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckCcnBandwidthLimitExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		ccnID := rs.Primary.Attributes["ccn_id"]
		info, has, err := service.DescribeCcn(ctx, ccnID)
		if err != nil {
			return err
		}
		if has == 0 {
			return fmt.Errorf("ccn instance not exists")
		}
		bandwidth, err := service.GetCcnRegionBandwidthLimit(ctx,
			ccnID,
			rs.Primary.Attributes["region"],
			rs.Primary.Attributes["dst_region"],
			info.bandWithLimitType)

		if err != nil {
			return err
		}

		if fmt.Sprintf("%d", bandwidth) != rs.Primary.Attributes["bandwidth_limit"] {
			return fmt.Errorf("ccn attachment not exists")
		}
		return nil
	}
}

const ccnOuterBase = `
resource tencentcloud_ccn main {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}`

const testAccCcnOuterBandwidthLimitConfig = ccnOuterBase + `
variable "other_region1" {
  default = "ap-shanghai"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  bandwidth_limit = 500
}
`

const testAccCcnOuterBandwidthLimitConfigUpdate = ccnOuterBase + `
variable "other_region1" {
  default = "ap-shanghai"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  bandwidth_limit = 100
}
`

const ccnInterBase = `
resource tencentcloud_ccn main {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
  qos                  = "AG"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
}
`

const testAccCcnInterBandwidthLimitConfig = ccnInterBase + `
variable "other_region1" {
  default = "ap-shanghai"
}

variable "other_region2" {
  default = "ap-beijing"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  dst_region      = var.other_region2
  bandwidth_limit = 500
}
`

const testAccCcnInterBandwidthLimitConfigUpdate = ccnInterBase + `
variable "other_region1" {
  default = "ap-shanghai"
}

variable "other_region2" {
  default = "ap-nanjing"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  dst_region      = var.other_region2
  bandwidth_limit = 100
}
`
