package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafCcSessionResource_basic -v
func TestAccTencentCloudWafCcSessionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCcSession,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_cc_session.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "domain", "keep.qcloudwaf.com"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "source", "get"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "category", "match"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "key_or_start_mat", "key_a=123"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "end_mat", "&"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "start_offset", "-1"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "end_offset", "-1"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "edition", "sparta-waf"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "session_name", "terraformDemo"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_cc_session.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafCcSessionUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_cc_session.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "domain", "keep.qcloudwaf.com"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "source", "post"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "category", "match"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "key_or_start_mat", "key_a=456"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "end_mat", "&"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "start_offset", "-1"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "end_offset", "-1"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "edition", "sparta-waf"),
					resource.TestCheckResourceAttr("tencentcloud_waf_cc_session.example", "session_name", "terraformDemo"),
				),
			},
		},
	})
}

const testAccWafCcSession = `
resource "tencentcloud_waf_cc_session" "example" {
  domain           = "keep.qcloudwaf.com"
  source           = "get"
  category         = "match"
  key_or_start_mat = "key_a=123"
  end_mat          = "&"
  start_offset     = "-1"
  end_offset       = "-1"
  edition          = "sparta-waf"
  session_name     = "terraformDemo"
}
`

const testAccWafCcSessionUpdate = `
resource "tencentcloud_waf_cc_session" "example" {
  domain           = "keep.qcloudwaf.com"
  source           = "post"
  category         = "match"
  key_or_start_mat = "key_a=456"
  end_mat          = "&"
  start_offset     = "-1"
  end_offset       = "-1"
  edition          = "sparta-waf"
  session_name     = "terraformDemo"
}
`
