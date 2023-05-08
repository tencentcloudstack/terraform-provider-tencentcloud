package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonitorPolicyConditions(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitorPolicyConditions(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_policy_conditions.monitor_policy_conditions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_policy_conditions.monitor_policy_conditions",
						"list.#"),
				),
			},
		},
	})
}

func testAccDataSourceMonitorPolicyConditions() string {
	return `
data "tencentcloud_monitor_policy_conditions" "monitor_policy_conditions" {
  name               = "Cloud Virtual Machine"
}`
}
