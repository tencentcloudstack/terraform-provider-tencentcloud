package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCssLiveTranscodeRuleAttachment_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssLiveTranscodeRuleAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_rule_attachment.live_transcode_rule_attachment", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_live_transcode_rule_attachment.liveTranscodeRuleAttachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssLiveTranscodeRuleAttachment = `

resource "tencentcloud_css_live_transcode_rule_attachment" "live_transcode_rule_attachment" {
  domain_name = ""
  app_name = ""
  stream_name = ""
  template_id = ""
    }

`
