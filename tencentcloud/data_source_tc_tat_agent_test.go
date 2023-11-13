package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tat_agent.agent")),
			},
		},
	})
}

const testAccTatAgentDataSource = `

data "tencentcloud_tat_agent" "agent" {
  instance_ids = 
  filters {
		name = ""
		values = 

  }
  limit = 
  offset = 
  }

`
