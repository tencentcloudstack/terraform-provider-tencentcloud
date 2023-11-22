package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudWafCcSessionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCcSession,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_waf_cc_session.cc_session", "id")),
			},
			{
				ResourceName:      "tencentcloud_waf_cc_session.cc_session",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafCcSession = `

resource "tencentcloud_waf_cc_session" "cc_session" {
  domain = "www.testwaf.com"
  source = "get"
  category = "match"
  key_or_start_mat = "key_a="
  end_mat = "&amp;"
  start_offset = "1"
  end_offset = "12"
  edition = "clb-waf"
  session_name = "name1"
  session_i_d = 122211
}

`
