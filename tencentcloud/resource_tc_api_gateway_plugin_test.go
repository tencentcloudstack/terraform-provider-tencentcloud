package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApiGatewayPluginResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayPlugin,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin.plugin", "id")),
			},
			{
				ResourceName:      "tencentcloud_api_gateway_plugin.plugin",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApiGatewayPlugin = `

resource "tencentcloud_api_gateway_plugin" "plugin" {
  plugin_name = ""
  plugin_type = ""
  plugin_data = ""
  description = ""
}

`
