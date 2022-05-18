package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudPrivateDnsZone_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsZone_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_private_dns_zone.zone", "domain", "domain.com"),
				),
			},
			{
				ResourceName:      "tencentcloud_private_dns_zone.zone",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPrivateDnsZone_basic = defaultInstanceVariable + `
resource "tencentcloud_private_dns_zone" "zone" {
  dns_forward_status = "DISABLED"
  domain             = "domain.com"
  remark             = "test_zone"
  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = var.cvm_vpc_id
  }
  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = var.vpc_id
  }
  tags = {
    "created-by" : "terraform",
  }
}
`
