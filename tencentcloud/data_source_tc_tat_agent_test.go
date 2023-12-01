package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTatAgentDataSource_basic -v
func TestAccTencentCloudTatAgentDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatAgentDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tat_agent.agent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_agent.agent", "automation_agent_set.0.agent_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_agent.agent", "automation_agent_set.0.environment"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_agent.agent", "automation_agent_set.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_agent.agent", "automation_agent_set.0.last_heartbeat_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_agent.agent", "automation_agent_set.0.support_features.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tat_agent.agent", "automation_agent_set.0.version"),
				),
			},
		},
	})
}

const testAccTatAgentDataSource = `

data "tencentcloud_tat_agent" "agent" {
	# instance_ids = ["ins-f9jr4bd2"]
	filters {
		  name = "environment"
		  values = ["Linux"]
	}
}

`
