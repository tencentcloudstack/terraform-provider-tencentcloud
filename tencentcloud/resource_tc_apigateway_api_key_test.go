package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApigatewayApiKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayApiKey,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_apigateway_api_key.api_key", "id")),
			},
			{
				ResourceName:      "tencentcloud_apigateway_api_key.api_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApigatewayApiKey = `

resource "tencentcloud_apigateway_api_key" "api_key" {
  secret_name = ""
  access_key_type = ""
  access_key_id = ""
  access_key_secret = ""
}

`
