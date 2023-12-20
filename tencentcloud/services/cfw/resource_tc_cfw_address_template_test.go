package cfw_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwAddressTemplateResource_basic -v
func TestAccTencentCloudNeedFixCfwAddressTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwAddressTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_address_template.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_address_template.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_address_template.example", "detail"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_address_template.example", "ip_string"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_address_template.example", "type"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_address_template.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCfwAddressTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_address_template.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_address_template.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_address_template.example", "detail"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_address_template.example", "ip_string"),
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_address_template.example", "type"),
				),
			},
		},
	})
}

const testAccCfwAddressTemplate = `
resource "tencentcloud_cfw_address_template" "example" {
  name      = "tf_example"
  detail    = "test template"
  ip_string = "1.1.1.1,2.2.2.2"
  type      = 1
}
`

const testAccCfwAddressTemplateUpdate = `
resource "tencentcloud_cfw_address_template" "example" {
  name      = "tf_example_update"
  detail    = "test template update"
  ip_string = "www.qq.com,www.tencent.com"
  type      = 5
}
`
