package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCcnV3BandwidthLimitBasic(t *testing.T) {

	keyNameLimit1 := "tencentcloud_ccn_bandwidth_limit.limit1"
	keyNameLimit2 := "tencentcloud_ccn_bandwidth_limit.limit2"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnBandwidthLimitConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnBandwidthLimitExists(keyNameLimit1),
					testAccCheckCcnBandwidthLimitExists(keyNameLimit2),
					resource.TestCheckResourceAttrSet(keyNameLimit1, "ccn_id"),
					resource.TestCheckResourceAttrSet(keyNameLimit1, "region"),
					resource.TestCheckResourceAttr(keyNameLimit1, "bandwidth_limit", "500"),

					resource.TestCheckResourceAttrSet(keyNameLimit2, "ccn_id"),
					resource.TestCheckResourceAttrSet(keyNameLimit2, "region"),
					resource.TestCheckResourceAttrSet(keyNameLimit2, "bandwidth_limit"),
				),
			},
		},
	})
}

func testAccCheckCcnBandwidthLimitExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := GetLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		bandwidth, err := service.DescribeCcnRegionBandwidthLimit(ctx,
			rs.Primary.Attributes["ccn_id"],
			rs.Primary.Attributes["region"])

		if err != nil {
			return err
		}

		if fmt.Sprintf("%d", bandwidth) != rs.Primary.Attributes["bandwidth_limit"] {
			return fmt.Errorf("ccn attachment not exists.")
		}
		return nil
	}
}

const testAccCcnBandwidthLimitConfig = `
variable "other_region1" {
    default = "ap-shanghai"
}

variable "other_region2" {
    default = "ap-beijing"
}

resource tencentcloud_ccn main{
	name ="ci-temp-test-ccn"
	description="ci-temp-test-ccn-des"
	qos ="AG"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
	ccn_id ="${tencentcloud_ccn.main.id}"
	region ="${var.other_region1}"
	bandwidth_limit = 500
}

resource tencentcloud_ccn_bandwidth_limit limit2 {
	ccn_id ="${tencentcloud_ccn.main.id}"
	region ="${var.other_region2}"
}
`
