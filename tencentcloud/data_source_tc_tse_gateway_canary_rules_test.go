package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateway_canary_rules.gateway_canary_rules")),
			},
		},
	})
}

const testAccTseGatewayCanaryRulesDataSource = `

data "tencentcloud_tse_gateway_canary_rules" "gateway_canary_rules" {
  gateway_id = "gateway-xxxxxx"
  service_id = "451a9920-e67a-4519-af41-fccac0e72005"
  }

`
