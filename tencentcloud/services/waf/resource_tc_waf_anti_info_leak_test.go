package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafAntiInfoLeakResource_basic -v
func TestAccTencentCloudWafAntiInfoLeakResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafAntiInfoLeak,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_anti_info_leak.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_anti_info_leak.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafAntiInfoLeakUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_anti_info_leak.example", "id"),
				),
			},
		},
	})
}

const testAccWafAntiInfoLeak = `
resource "tencentcloud_waf_anti_info_leak" "example" {
  domain      = "keep.qcloudwaf.com"
  name        = "tf_example"
  action_type = 0
  strategies {
    field   = "information"
    content = "phone"
  }
  uri    = "/anti_info_leak_url"
  status = 0
}
`

const testAccWafAntiInfoLeakUpdate = `
resource "tencentcloud_waf_anti_info_leak" "example" {
  domain      = "keep.qcloudwaf.com"
  name        = "tf_example_update"
  action_type = 0
  strategies {
    field   = "returncode"
    content = "400"
  }
  uri    = "/anti_info_leak_url"
  status = 1
}
`
