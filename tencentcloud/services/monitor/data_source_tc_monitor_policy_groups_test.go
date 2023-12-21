package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonitorPolicyGroups(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMonitorPolicyGroups(),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_policy_groups.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_policy_groups.name",
						"list.#"),
				),
			},
		},
	})
}

func testAccDataSourceMonitorPolicyGroups() string {
	return `data "tencentcloud_monitor_policy_groups" "name" {}`
}
