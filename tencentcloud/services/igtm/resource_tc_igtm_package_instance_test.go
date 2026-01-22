package igtm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIgtmPackageInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIgtmPackage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_package_instance.igtm_package", "id"),
				),
			},
			{
				Config: testAccIgtmPackageUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_package_instance.igtm_package", "id"),
				),
			},
		},
	})
}

const testAccIgtmPackage = `
resource "tencentcloud_igtm_package_instance" "example" {
  goods_type   = "STANDARD"
  auto_renew   = 1
  time_span    = 1
  auto_voucher = 1
}
`

const testAccIgtmPackageUpdate = `
resource "tencentcloud_igtm_package_instance" "example" {
  goods_type   = "ULTIMATE"
  auto_renew   = 2
  time_span    = 2
  auto_voucher = 0
}
`
