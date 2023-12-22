package tse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTseGatewaysDataSource_basic -v
func TestAccTencentCloudTseGatewaysDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGatewaysDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateways.gateways"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.auto_renew_flag"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.cur_deadline"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.enable_cls"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.enable_internet"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.engine_region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.feature_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.gateway_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.gateway_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.ingress_class_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.instance_port.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.internet_max_bandwidth_out"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.internet_pay_mode"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.isolate_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.node_config.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.node_config.0.number"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.node_config.0.specification"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.trade_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.vpc_config.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.vpc_config.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.gateway_list.0.vpc_config.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateways.gateways", "result.0.total_count"),
				),
			},
		},
	})
}

const testAccTseGatewaysDataSource = `

data "tencentcloud_tse_gateways" "gateways" {
  filters {
    name   = "GatewayId"
    values = ["gateway-ddbb709b"]
  }
}

`
