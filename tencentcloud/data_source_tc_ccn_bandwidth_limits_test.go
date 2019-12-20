package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudCcnV3BandwidthLimitsBasic(t *testing.T) {
	keyName := "data.tencentcloud_ccn_bandwidth_limits.limit"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudCcnBandwidthLimits,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(keyName),
					resource.TestCheckResourceAttrSet(keyName, "ccn_id"),
					resource.TestCheckResourceAttrSet(keyName, "limits.#"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudCcnBandwidthLimits = `

variable "other_region1" {
  default = "ap-shanghai"
}

resource tencentcloud_ccn main {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

data tencentcloud_ccn_bandwidth_limits limit {
  ccn_id = "${tencentcloud_ccn.main.id}"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
  ccn_id          = "${tencentcloud_ccn.main.id}"
  region          = "${var.other_region1}"
  bandwidth_limit = 500
}
`
