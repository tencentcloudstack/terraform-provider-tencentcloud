package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSesTemplate_basic -v
func TestAccTencentCloudSesTemplate_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckBusiness(t, ACCOUNT_TYPE_SES) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ses_template.template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ses_template.template", "template_name", "sesTemplateName"),
				),
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
  template_name = "sesTemplateName"
  template_content {
    text = "This is the content of the test"
  }
}

`
