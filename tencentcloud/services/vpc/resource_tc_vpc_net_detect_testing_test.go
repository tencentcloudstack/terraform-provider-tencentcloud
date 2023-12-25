package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTestingVpcNetDetectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingVpcNetDetect,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_net_detect.net_detect", "id")),
			},
			{
				Config: testAccTestingVpcNetDetectUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_net_detect.net_detect", "net_detect_name", "terraform-for-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_net_detect.net_detect",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTestingVpcNetDetect = `

resource "tencentcloud_vpc_net_detect" "net_detect" {
  net_detect_name       = "terraform-test"
  vpc_id                = "vpc-jxnxbc07"
  subnet_id             = "subnet-ev908x0w"
  next_hop_destination  = "nat-bfnnl8wg"
  next_hop_type         = "NAT"
  detect_destination_ip = [
    "172.16.128.110"
  ]
}

`

const testAccTestingVpcNetDetectUpdate = `

resource "tencentcloud_vpc_net_detect" "net_detect" {
  net_detect_name       = "terraform-for-test"
  vpc_id                = "vpc-jxnxbc07"
  subnet_id             = "subnet-ev908x0w"
  next_hop_destination  = "nat-bfnnl8wg"
  next_hop_type         = "NAT"
  detect_destination_ip = [
    "172.16.128.110"
  ]
}

`
