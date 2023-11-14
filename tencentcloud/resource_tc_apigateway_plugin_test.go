package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApigatewayPluginResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayPlugin,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_apigateway_plugin.plugin", "id")),
			},
			{
				ResourceName:      "tencentcloud_apigateway_plugin.plugin",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApigatewayPlugin = `

resource "tencentcloud_apigateway_plugin" "plugin" {
  plugin_name = ""
  plugin_type = ""
  plugin_data = ""
  description = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
