package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssTimeshiftRuleAttachmentResource_basic -v
func TestAccTencentCloudCssTimeshiftRuleAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssTimeshiftRuleAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_timeshift_rule_attachment.timeshift_rule_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_rule_attachment.timeshift_rule_attachment", "domain_name", "177154.push.tlivecloud.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_timeshift_rule_attachment.timeshift_rule_attachment", "template_id"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_rule_attachment.timeshift_rule_attachment", "app_name", "qqq"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_rule_attachment.timeshift_rule_attachment", "stream_name", "ppp"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_timeshift_rule_attachment.timeshift_rule_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssTimeshiftRuleAttachment = testAccCssTimeshiftTemplate + `

resource "tencentcloud_css_timeshift_rule_attachment" "timeshift_rule_attachment" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = tencentcloud_css_timeshift_template.timeshift_template.id
  app_name    = "qqq"
  stream_name = "ppp"
}

`
