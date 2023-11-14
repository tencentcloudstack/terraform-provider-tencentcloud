package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCssWatermarkRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssWatermarkRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_watermark_rule.watermark_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_css_watermark_rule.watermark_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssWatermarkRule = `

resource "tencentcloud_css_watermark_rule" "watermark_rule" {
  domain_name = &lt;nil&gt;
  app_name = &lt;nil&gt;
  stream_name = &lt;nil&gt;
  watermark_id = &lt;nil&gt;
    }

`
