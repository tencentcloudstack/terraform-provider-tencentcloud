package apigateway_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudAPIGatewayUpdateServiceResource_basic -v
func TestAccTencentCloudAPIGatewayUpdateServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIGatewayUpdateService1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_update_service.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_update_service.example", "service_id", "service-oczq2nyk"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_update_service.example", "environment_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_update_service.example", "version_name", "20240204142759-b5a4f741-adc0-4964-b01b-2a4a04ff6964"),
				),
			},
			{
				Config: testAccAPIGatewayUpdateService2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_update_service.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_update_service.example", "service_id", "service-oczq2nyk"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_update_service.example", "environment_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_update_service.example", "version_name", "20240126164018-c6ec85b5-8aae-4896-bc26-8afbf88dcbcc"),
				),
			},
		},
	})
}

const testAccAPIGatewayUpdateService1 = `
resource "tencentcloud_api_gateway_update_service" "example" {
  service_id       = "service-oczq2nyk"
  environment_name = "test"
  version_name     = "20240204142759-b5a4f741-adc0-4964-b01b-2a4a04ff6964"
}
`

const testAccAPIGatewayUpdateService2 = `
resource "tencentcloud_api_gateway_update_service" "example" {
  service_id       = "service-oczq2nyk"
  environment_name = "test"
  version_name     = "20240126164018-c6ec85b5-8aae-4896-bc26-8afbf88dcbcc"
}
`
