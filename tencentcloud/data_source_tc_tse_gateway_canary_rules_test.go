package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTseGatewayCanaryRulesDataSource_basic -v
func TestAccTencentCloudTseGatewayCanaryRulesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGatewayCanaryRulesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.balanced_service_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.balanced_service_list.0.percent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.balanced_service_list.0.service_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.balanced_service_list.0.service_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.balanced_service_list.0.upstream_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.condition_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.condition_list.0.key"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.condition_list.0.operator"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.condition_list.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.condition_list.0.value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.enabled"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.canary_rule_list.0.priority"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules", "result.0.total_count"),
				),
			},
		},
	})
}

const testAccTseGatewayCanaryRulesDataSource = `

data "tencentcloud_tse_gateway_canary_rules" "gateway_canary_rules" {
	gateway_id = "gateway-ddbb709b"
	service_id = "b6017eaf-2363-481e-9e93-8d65aaf498cd"
}

`
