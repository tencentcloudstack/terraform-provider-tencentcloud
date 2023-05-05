package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixApiGatewayPluginAttachmentResource_basic(t *testing.T) {
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
  plugin_id        = "plugin-ny74siyz"
  service_id       = "service-n1mgl0sq"
  environment_name = "test"
  api_id           = "api-6tfrdysk"
}

`
