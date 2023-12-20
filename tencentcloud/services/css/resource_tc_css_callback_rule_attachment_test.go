package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssCallbackRuleAttachmentResource_basic -v
func TestAccTencentCloudCssCallbackRuleAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssCallbackRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_callback_rule_attachment.callback_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_rule_attachment.callback_rule", "domain_name", "177154.push.tlivecloud.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_callback_rule_attachment.callback_rule", "template_id"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_rule_attachment.callback_rule", "app_name", "live"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_callback_rule_attachment.callback_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssCallbackRule = testAccCssCallbackTemplate + `

resource "tencentcloud_css_callback_rule_attachment" "callback_rule" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = tencentcloud_css_callback_template.callback_template.id
  app_name    = "live"
}

`
