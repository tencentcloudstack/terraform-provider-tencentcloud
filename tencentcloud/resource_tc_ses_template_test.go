package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ses_template.template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ses_template.template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSesTemplate = `

resource "tencentcloud_ses_template" "template" {
  template_name = "smsTemplateName"
  template_content {
		html = &lt;nil&gt;
		text = &lt;nil&gt;

  }
}

`
