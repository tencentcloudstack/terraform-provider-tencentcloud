package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssRecordRuleResource_basic -v
func TestAccTencentCloudCssRecordRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssRecordRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_record_rule_attachment.record_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_record_rule_attachment.record_rule", "domain_name", "177154.push.tlivecloud.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_record_rule_attachment.record_rule", "template_id"),
					resource.TestCheckResourceAttr("tencentcloud_css_record_rule_attachment.record_rule", "app_name", "qqq"),
					resource.TestCheckResourceAttr("tencentcloud_css_record_rule_attachment.record_rule", "stream_name", "ppp"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_record_rule_attachment.record_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssRecordRule = testAccCssRecordTemplate + `

resource "tencentcloud_css_record_rule_attachment" "record_rule" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = tencentcloud_css_record_template.record_template.id
  app_name    = "qqq"
  stream_name = "ppp"
}

`
