package ccn_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudCcnRouteTableInputPoliciesDataSource_basic -v
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
				resource.TestCheckResourceAttr("data.tencentcloud_ccn_route_table_input_policies.example", "policy_set.#", "1"),
			),
		}},
	})
}

const testAccCcnRouteTableInputPoliciesDataSource = `
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "description."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  tags = {
    createBy = "terraform"
  }
}

resource "tencentcloud_ccn_route_table" "example" {
  ccn_id      = tencentcloud_ccn.example.id
  name        = "tf-example"
  description = "desc."
}

resource "tencentcloud_ccn_route_table_input_policies" "example" {
  ccn_id         = tencentcloud_ccn.example.id
  route_table_id = tencentcloud_ccn_route_table.example.id
  policies {
    action      = "accept"
    description = "desc."
    route_conditions {
      name          = "instance-region"
      values        = ["ap-guangzhou"]
      match_pattern = 1
    }
  }
}

data "tencentcloud_ccn_route_table_input_policies" "example" {
  depends_on = [ tencentcloud_ccn_route_table_input_policies.example ]
  ccn_id         = tencentcloud_ccn.example.id
  route_table_id = tencentcloud_ccn_route_table.example.id
}
`
