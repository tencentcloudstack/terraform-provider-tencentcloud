package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafProtectionModeResource_basic -v
func TestAccTencentCloudWafProtectionModeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafProtectionMode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_protection_mode.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_protection_mode.example", "domain", "keep.qcloudwaf.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_protection_mode.example", "mode"),
					resource.TestCheckResourceAttr("tencentcloud_waf_protection_mode.example", "edition", "sparta-waf"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_protection_mode.example", "type"),
				),
			},
		},
	})
}

const testAccWafProtectionMode = `
resource "tencentcloud_waf_protection_mode" "example" {
  domain  = "keep.qcloudwaf.com"
  mode    = 10
  edition = "sparta-waf"
  type    = 0
}
`
