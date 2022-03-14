package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataAddressTemplatesNameAll = "data.tencentcloud_address_templates.all_test"

func TestAccTencentCloudDataAddressTemplates(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAddressTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataAddressTemplatesBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAddressTemplateExists("tencentcloud_address_template.myaddress_template"),
					resource.TestCheckResourceAttrSet(testDataAddressTemplatesNameAll, "template_list.0.id"),
					resource.TestCheckResourceAttr(testDataAddressTemplatesNameAll, "template_list.0.name", "testacctcrtemplate"),
					resource.TestCheckResourceAttr(testDataAddressTemplatesNameAll, "template_list.0.addresses.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataAddressTemplatesBasic = `
resource "tencentcloud_address_template" "myaddress_template" {
  name        = "testacctcrtemplate"
  addresses = ["1.1.1.1"]
}

data "tencentcloud_address_templates" "all_test" {
  name = tencentcloud_address_template.myaddress_template.name
  id = tencentcloud_address_template.myaddress_template.id
}

`
