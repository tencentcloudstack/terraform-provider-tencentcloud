package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
  plugin_name = "terraform-plugin-test"
  plugin_type = "IPControl"
  plugin_data = jsonencode({
    "type" : "white_list",
    "blocks" : "1.1.1.1\n2.2.2.2",
  })
  description = "terraform test"
}

`
