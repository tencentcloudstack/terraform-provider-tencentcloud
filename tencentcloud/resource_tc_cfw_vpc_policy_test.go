package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCfwVpcPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwVpcPolicy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_policy.vpc_policy", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfw_vpc_policy.vpc_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfwVpcPolicy = `

resource "tencentcloud_cfw_vpc_policy" "vpc_policy" {
  rules {
		source_content = "0.0.0.0/0"
		source_type = "net"
		dest_content = "192.168.0.2"
		dest_type = "net"
		protocol = "ANY"
		rule_action = "log"
		port = "-1/-1"
		description = "test vpc rule"
		order_index = 28
		uuid = 
		enable = "true"
		edge_id = "ALL"
		detected_times = 0
		edge_name = ""
		internal_uuid = 0
		deleted = 0
		fw_group_id = ""
		fw_group_name = ""
		beta_list {
			task_id = 
			task_name = ""
			last_time = ""
		}

  }
  from = ""
}

`
