package privatedns_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudPrivateDnsAccountResource_basic -v
func TestAccTencentCloudPrivateDnsAccountResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsAccount_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_private_dns_account.example", "account_uin", "100123456789"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_account.example", "account"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_account.example", "nickname"),
				),
			},
			{
				ResourceName:      "tencentcloud_private_dns_account.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPrivateDnsAccount_basic = `
resource "tencentcloud_private_dns_account" "example" {
  account_uin = "100123456789"
}
`
