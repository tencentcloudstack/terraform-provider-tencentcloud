package apigateway_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudApiGatewayApiAppApiDataSource_basic -v
func TestAccTencentCloudApiGatewayApiAppApiDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayApiAppApiDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_api_gateway_api_app_api.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_api_gateway_api_app_api.example", "service_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_api_gateway_api_app_api.example", "api_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_api_gateway_api_app_api.example", "api_region"),
				),
			},
		},
	})
}

const testAccApiGatewayApiAppApiDataSource = `
data "tencentcloud_api_gateway_api_app_api" "example" {
  service_id = "service-nxz6yync"
  api_id     = "api-0cvmf4x4"
  api_region = "ap-guangzhou"
}
`
