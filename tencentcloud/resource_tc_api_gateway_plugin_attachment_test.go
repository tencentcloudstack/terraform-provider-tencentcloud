package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApiGatewayPluginAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayPluginAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin_attachment.plugin_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_api_gateway_plugin_attachment.plugin_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApiGatewayPluginAttachment = `

resource "tencentcloud_api_gateway_plugin_attachment" "plugin_attachment" {
  plugin_id = ""
  service_id = ""
  environment_name = ""
  api_id = ""
}

`
