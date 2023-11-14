package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApigatewayApiAppResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayApiApp,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_apigateway_api_app.api_app", "id")),
			},
			{
				ResourceName:      "tencentcloud_apigateway_api_app.api_app",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApigatewayApiApp = `

resource "tencentcloud_apigateway_api_app" "api_app" {
  api_app_name = ""
  api_app_desc = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
