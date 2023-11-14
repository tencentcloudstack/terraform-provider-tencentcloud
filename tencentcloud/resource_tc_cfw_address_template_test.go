package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCfwAddressTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwAddressTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfw_address_template.address_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfw_address_template.address_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfwAddressTemplate = `

resource "tencentcloud_cfw_address_template" "address_template" {
  name = "test address template"
  detail = "test address template"
  ip_string = "1.1.1.1,2.2.2.2"
  type = 1
}

`
