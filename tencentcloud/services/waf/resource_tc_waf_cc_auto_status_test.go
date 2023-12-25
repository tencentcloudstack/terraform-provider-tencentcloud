package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafCcAutoStatusResource_basic -v
func TestAccTencentCloudWafCcAutoStatusResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCcAutoStatus,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_cc_auto_status.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_auto_status.example", "domain", "keep.qcloudwaf.com"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_auto_status.example", "edition", "sparta-waf"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_cc_auto_status.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafCcAutoStatus = `
resource "tencentcloud_waf_cc_auto_status" "example" {
  domain  = "keep.qcloudwaf.com"
  edition = "sparta-waf"
}
`
