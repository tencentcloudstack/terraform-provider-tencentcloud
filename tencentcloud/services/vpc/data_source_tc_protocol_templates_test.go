package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataProtocolTemplatesNameAll = "data.tencentcloud_protocol_templates.all_test"

func TestAccTencentCloudDataProtocolTemplates(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckProtocolTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataProtocolTemplatesBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProtocolTemplateExists("tencentcloud_protocol_template.myprotocol_template"),
					resource.TestCheckResourceAttrSet(testDataProtocolTemplatesNameAll, "template_list.0.id"),
					resource.TestCheckResourceAttr(testDataProtocolTemplatesNameAll, "template_list.0.name", "testacctcrtemplate"),
					resource.TestCheckResourceAttr(testDataProtocolTemplatesNameAll, "template_list.0.protocols.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataProtocolTemplatesBasic = `
resource "tencentcloud_protocol_template" "myprotocol_template" {
  name        = "testacctcrtemplate"
  protocols = ["udp:all"]
}

data "tencentcloud_protocol_templates" "all_test" {
  name = tencentcloud_protocol_template.myprotocol_template.name
}

`
