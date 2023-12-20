package apigateway_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudApiGatewayServiceReleaseVersionsDataSource_basic -v
func TestAccTencentCloudApiGatewayServiceReleaseVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayServiceReleaseVersionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_api_gateway_service_release_versions.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_api_gateway_service_release_versions.example", "service_id"),
				),
			},
		},
	})
}

const testAccApiGatewayServiceReleaseVersionsDataSource = `
data "tencentcloud_api_gateway_service_release_versions" "example" {
  service_id = "service-nxz6yync"
}
`
