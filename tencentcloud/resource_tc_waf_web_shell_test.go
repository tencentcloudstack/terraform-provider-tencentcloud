package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafWebShellResource_basic -v
func TestAccTencentCloudWafWebShellResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafWebShell,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_web_shell.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_web_shell.example", "domain", "keep.qcloudwaf.com"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_web_shell.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafWebShellUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_web_shell.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_web_shell.example", "domain", "keep.qcloudwaf.com"),
				),
			},
		},
	})
}

const testAccWafWebShell = `
resource "tencentcloud_waf_web_shell" "example" {
  domain = "keep.qcloudwaf.com"
  status = 1
}
`

const testAccWafWebShellUpdate = `
resource "tencentcloud_waf_web_shell" "example" {
  domain = "keep.qcloudwaf.com"
  status = 0
}
`
