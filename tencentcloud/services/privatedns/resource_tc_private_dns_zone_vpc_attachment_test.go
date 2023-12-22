package privatedns_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudPrivateDnsZoneVpcAttachmentResource_basic -v
func TestAccTencentCloudPrivateDnsZoneVpcAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsZoneVpcAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_zone_vpc_attachment.example1", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_private_dns_zone_vpc_attachment.example1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPrivateDnsZoneAccountVpcAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_zone_vpc_attachment.example2", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_private_dns_zone_vpc_attachment.example2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPrivateDnsZoneVpcAttachment = `
resource "tencentcloud_private_dns_zone_vpc_attachment" "example1" {
  zone_id = "zone-980faacc"

  vpc_set {
    uniq_vpc_id = "vpc-86v957zb"
    region      = "ap-guangzhou"
  }
}
`

const testAccPrivateDnsZoneAccountVpcAttachment = `
resource "tencentcloud_private_dns_zone_vpc_attachment" "example2" {
  zone_id = "zone-980faacc"

  account_vpc_set {
    uniq_vpc_id = "vpc-axrsmmrv"
    region      = "ap-guangzhou"
    uin         = "100022770164"
  }
}
`
