package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonitorPolicyConditions(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitorPolicyConditions(),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_policy_conditions.monitor_policy_conditions"),
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
