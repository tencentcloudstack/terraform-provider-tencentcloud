package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCssLiveTranscodeRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssLiveTranscodeRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_rule.live_transcode_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_css_live_transcode_rule.live_transcode_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssLiveTranscodeRule = `

resource "tencentcloud_css_live_transcode_rule" "live_transcode_rule" {
  domain_name = &lt;nil&gt;
  app_name = ""
  stream_name = ""
  template_id = &lt;nil&gt;
    }

`
