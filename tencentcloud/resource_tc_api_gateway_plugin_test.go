package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudApiGatewayPluginResource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin.example", "plugin_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin.example", "plugin_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin.example", "plugin_data"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin.example", "description"),
				),
			},
			{
				ResourceName:      "tencentcloud_api_gateway_plugin.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApiGatewayPluginUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin.example", "plugin_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin.example", "plugin_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin.example", "plugin_data"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin.example", "description"),
				),
			},
		},
	})
}

const testAccApiGatewayPlugin = `
resource "tencentcloud_api_gateway_plugin" "example" {
  plugin_name = "tf-example"
  plugin_type = "IPControl"
  plugin_data = jsonencode({
    "type" : "white_list",
    "blocks" : "1.1.1.1",
  })
  description = "desc."
}
`

const testAccApiGatewayPluginUpdate = `
resource "tencentcloud_api_gateway_plugin" "example" {
  plugin_name = "tf-example-update"
  plugin_type = "IPControl"
  plugin_data = jsonencode({
    "type" : "white_list",
    "blocks" : "2.2.2.2",
  })
  description = "update desc."
}
`
