package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataServiceTemplatesNameAll = "data.tencentcloud_service_templates.all_test"

func TestAccTencentCloudDataServiceTemplates(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataServiceTemplatesBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceTemplateExists("tencentcloud_service_template.myservice_template"),
					resource.TestCheckResourceAttrSet(testDataServiceTemplatesNameAll, "template_list.0.id"),
					resource.TestCheckResourceAttr(testDataServiceTemplatesNameAll, "template_list.0.name", "testacctcrtemplate"),
					resource.TestCheckResourceAttr(testDataServiceTemplatesNameAll, "template_list.0.services.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataServiceTemplatesBasic = `
resource "tencentcloud_service_template" "myservice_template" {
  name        = "testacctcrtemplate"
  services = ["udp:all"]
}

data "tencentcloud_service_templates" "all_test" {
  name = tencentcloud_service_template.myservice_template.name
}

`
