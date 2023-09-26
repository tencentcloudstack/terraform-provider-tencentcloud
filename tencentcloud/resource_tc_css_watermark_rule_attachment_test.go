package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCssWatermarkRuleAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssWatermarkRuleAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_css_watermark_rule_attachment.watermark_rule_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssWatermarkRuleAttachment = `

resource "tencentcloud_css_watermark_rule_attachment" "watermark_rule_attachment" {
  domain_name = &lt;nil&gt;
  app_name = &lt;nil&gt;
  stream_name = &lt;nil&gt;
  template_id = 
    }

`
