package ccn_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCcnRouteTableInputPoliciesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCcnRouteTableInputPoliciesDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ccn_route_table_input_policies.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_route_table_input_policies.example", "ccn_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_ccn_route_table_input_policies.example", "route_table_id"),
			),
		}},
	})
}

const testAccCcnRouteTableInputPoliciesDataSource = `
data "tencentcloud_ccn_route_table_input_policies" "example" {
  ccn_id         = "ccn-06jek8tf"
  route_table_id = "ccnrtb-4jv5ltb9"
}
`
