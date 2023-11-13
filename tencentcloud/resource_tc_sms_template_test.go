package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSmsTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSmsTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sms_template.template", "id")),
			},
			{
				ResourceName:      "tencentcloud_sms_template.template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSmsTemplate = `

resource "tencentcloud_sms_template" "template" {
  template_name = "TemplateName"
  template_content = "Template test content"
  international = 0
  sms_type = 0
  remark = "sms test"
      }

`
