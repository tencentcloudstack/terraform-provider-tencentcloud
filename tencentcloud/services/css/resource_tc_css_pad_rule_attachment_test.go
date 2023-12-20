package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssPadRuleAttachmentResource_basic -v
func TestAccTencentCloudCssPadRuleAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssPadRuleAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_pad_rule_attachment.pad_rule_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_rule_attachment.pad_rule_attachment", "domain_name", "177154.push.tlivecloud.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pad_rule_attachment.pad_rule_attachment", "template_id"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_rule_attachment.pad_rule_attachment", "app_name", "qqq"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_rule_attachment.pad_rule_attachment", "stream_name", "ppp"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_pad_rule_attachment.pad_rule_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssPadRuleAttachment = testAccCssPadTemplate + `

resource "tencentcloud_css_pad_rule_attachment" "pad_rule_attachment" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = tencentcloud_css_pad_template.pad_template.id
  app_name    = "qqq"
  stream_name = "ppp"
}

`
