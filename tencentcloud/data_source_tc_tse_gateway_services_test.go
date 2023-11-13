package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseGatewayServicesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGatewayServicesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateway_services.gateway_services")),
			},
		},
	})
}

const testAccTseGatewayServicesDataSource = `

data "tencentcloud_tse_gateway_services" "gateway_services" {
  gateway_id = "gateway-xxxxxx"
  filters {
		key = "name"
		value = "serviceA"

  }
  }

`
