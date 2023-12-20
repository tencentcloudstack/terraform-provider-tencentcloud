package apigateway_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudApiGatewayPluginAttachmentResource_basic -v
func TestAccTencentCloudApiGatewayPluginAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayPluginAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_plugin_attachment.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_api_gateway_plugin_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApiGatewayPluginAttachment = `
resource "tencentcloud_api_gateway_plugin" "example" {
  plugin_name = "tf-example"
  plugin_type = "IPControl"
  plugin_data = jsonencode({
    "type" : "white_list",
    "blocks" : "1.1.1.1",
  })
  description = "desc."
}

resource "tencentcloud_api_gateway_service" "example" {
  service_name = "tf_example_service"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_api" "example" {
  service_id            = tencentcloud_api_gateway_service.example.id
  api_name              = "tf_example_api"
  api_desc              = "desc."
  auth_type             = "APP"
  protocol              = "HTTP"
  enable_cors           = true
  request_config_path   = "/user/info"
  request_config_method = "GET"

  request_parameters {
    name          = "name"
    position      = "QUERY"
    type          = "string"
    desc          = "desc."
    default_value = "terraform"
    required      = true
  }

  service_config_type      = "HTTP"
  service_config_timeout   = 15
  service_config_url       = "https://www.qq.com"
  service_config_path      = "/user"
  service_config_method    = "GET"
  response_type            = "HTML"
  response_success_example = "success"
  response_fail_example    = "fail"

  response_error_codes {
    code           = 400
    msg            = "system error msg."
    desc           = "system error desc."
    converted_code = 407
    need_convert   = true
  }
}

resource "tencentcloud_api_gateway_service_release" "example" {
  service_id       = tencentcloud_api_gateway_api.example.service_id
  environment_name = "release"
  release_desc     = "desc."
}

resource "tencentcloud_api_gateway_plugin_attachment" "example" {
  plugin_id        = tencentcloud_api_gateway_plugin.example.id
  service_id       = tencentcloud_api_gateway_service_release.example.service_id
  api_id           = tencentcloud_api_gateway_api.example.id
  environment_name = "release"
}
`
