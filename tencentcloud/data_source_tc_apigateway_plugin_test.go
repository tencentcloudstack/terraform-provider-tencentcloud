package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApigatewayPluginDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayPluginDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_apigateway_plugin.plugin")),
			},
		},
	})
}

const testAccApigatewayPluginDataSource = `

data "tencentcloud_apigateway_plugin" "plugin" {
  service_id = ""
  plugin_id = ""
  environment_name = ""
    tags = {
    "createdBy" = "terraform"
  }
}

`
