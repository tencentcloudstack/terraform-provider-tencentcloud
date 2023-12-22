package tse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTseGatewayServicesDataSource_basic -v
func TestAccTencentCloudTseGatewayServicesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGatewayServicesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateway_services.gateway_services"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.created_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.editable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.tags.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.upstream_info.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.upstream_info.0.algorithm"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.upstream_info.0.host"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.upstream_info.0.targets.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.upstream_info.0.targets.0.created_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.upstream_info.0.targets.0.health"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.upstream_info.0.targets.0.host"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.upstream_info.0.targets.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.upstream_info.0.targets.0.weight"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.service_list.0.upstream_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_services.gateway_services", "result.0.total_count"),
				),
			},
		},
	})
}

const testAccTseGatewayServicesDataSource = `

data "tencentcloud_tse_gateway_services" "gateway_services" {
	gateway_id = "gateway-ddbb709b"
	filters {
	  key   = "name"
	  value = "test"
	}
}

`
