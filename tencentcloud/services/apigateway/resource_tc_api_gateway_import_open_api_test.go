package apigateway_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudApiGatewayImportOpenApiResource_basic -v
func TestAccTencentCloudApiGatewayImportOpenApiResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayImportOpenApi,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_import_open_api.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_import_open_api.example", "service_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_import_open_api.example", "content"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_import_open_api.example", "encode_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_import_open_api.example", "content_version"),
				),
			},
		},
	})
}

const testAccApiGatewayImportOpenApi = `
resource "tencentcloud_api_gateway_import_open_api" "example" {
  service_id      = "service-nxz6yync"
  content         = "info:\n  title: keep-service\n  version: 1.0.1\nopenapi: 3.0.0\npaths:\n  /api/test:\n    get:\n      description: desc\n      operationId: test\n      responses:\n        '200':\n          content:\n            text/html:\n              example: '200'\n          description: '200'\n        default:\n          content:\n            text/html:\n              example: '400'\n          description: '400'\n      x-apigw-api-business-type: NORMAL\n      x-apigw-api-type: NORMAL\n      x-apigw-backend:\n        ServiceConfig:\n          Method: GET\n          Path: /test\n          Url: http://domain.com\n        ServiceType: HTTP\n      x-apigw-cors: false\n      x-apigw-protocol: HTTP\n      x-apigw-service-timeout: 15\n"
  encode_type     = "YAML"
  content_version = "openAPI"
}
`
