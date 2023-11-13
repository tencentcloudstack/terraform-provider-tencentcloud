package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafModuleStatusResource_basic -v
func TestAccTencentCloudWafModuleStatusResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafModuleStatus,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_module_status.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_module_status.example", "domain", "keep.qcloudwaf.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_module_status.example", "web_security"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_module_status.example", "access_control"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_module_status.example", "cc_protection"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_module_status.example", "api_protection"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_module_status.example", "anti_tamper"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_module_status.example", "anti_leakage"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_module_status.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafModuleStatus = `
resource "tencentcloud_waf_module_status" "example" {
  domain         = "keep.qcloudwaf.com"
  web_security   = 1
  access_control = 0
  cc_protection  = 1
  api_protection = 1
  anti_tamper    = 1
  anti_leakage   = 0
}
`
