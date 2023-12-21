package privatedns_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudPrivateDnsZone_basic -v
func TestAccTencentCloudPrivateDnsZone_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsZone_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_private_dns_zone.example", "domain", "domain.com"),
				),
			},
			{
				ResourceName:      "tencentcloud_private_dns_zone.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPrivateDnsZone_basic = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_private_dns_zone" "example" {
  domain = "domain.com"
  remark = "remark."

  vpc_set {
    region      = "ap-guangzhou"
    uniq_vpc_id = var.cvm_vpc_id
  }

  dns_forward_status   = "DISABLED"
  cname_speedup_status = "ENABLED"

  tags = {
    createdBy : "terraform"
  }
}
`
