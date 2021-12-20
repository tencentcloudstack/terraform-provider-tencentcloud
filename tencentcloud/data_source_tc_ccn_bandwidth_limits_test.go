package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudCcnV3BandwidthLimitsOuter(t *testing.T) {
	t.Parallel()
	keyName := "data.tencentcloud_ccn_bandwidth_limits.limit"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudCcnOuterBandwidthLimits,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(keyName),
					resource.TestCheckResourceAttrSet(keyName, "ccn_id"),
					resource.TestCheckResourceAttr(keyName, "limits.#", "1"),
					resource.TestCheckResourceAttr(keyName, "limits.0.region", "ap-shanghai"),
					resource.TestCheckResourceAttr(keyName, "limits.0.bandwidth_limit", "500"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudCcnV3BandwidthLimitsInter(t *testing.T) {
	t.Parallel()
	keyName := "data.tencentcloud_ccn_bandwidth_limits.limit"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudCcnInterBandwidthLimits,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(keyName),
					resource.TestCheckResourceAttrSet(keyName, "ccn_id"),
					resource.TestCheckResourceAttr(keyName, "limits.#", "1"),
					resource.TestCheckResourceAttr(keyName, "limits.0.region", "ap-shanghai"),
					resource.TestCheckResourceAttr(keyName, "limits.0.dst_region", "ap-beijing"),
					resource.TestCheckResourceAttr(keyName, "limits.0.bandwidth_limit", "500"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudCcnOuterBandwidthLimits = `
variable "other_region1" {
  default = "ap-shanghai"
}

resource tencentcloud_ccn main {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  bandwidth_limit = 500
}

data tencentcloud_ccn_bandwidth_limits limit {
  ccn_id = tencentcloud_ccn_bandwidth_limit.limit1.ccn_id
}
`

const TestAccDataSourceTencentCloudCcnInterBandwidthLimits = `
variable "other_region1" {
  default = "ap-shanghai"
}

variable "other_region2" {
  default = "ap-beijing"
}

resource tencentcloud_ccn main {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
  qos                  = "AG"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  dst_region      = var.other_region2
  bandwidth_limit = 500
}

data tencentcloud_ccn_bandwidth_limits limit {
  ccn_id = tencentcloud_ccn_bandwidth_limit.limit1.ccn_id
}
`
