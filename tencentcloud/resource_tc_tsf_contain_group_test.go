package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTsfContainGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfContainGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_contain_group.contain_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_contain_group.contain_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfContainGroup = `

resource "tencentcloud_tsf_contain_group" "contain_group" {
  application_id = ""
  namespace_id = ""
  group_name = ""
  instance_num = 
  access_type = 
  protocol_ports {
		protocol = ""
		port = 
		target_port = 
		node_port = 

  }
  cluster_id = ""
  cpu_limit = ""
  mem_limit = ""
  group_comment = ""
  update_type = 
  update_ivl = 
  cpu_request = ""
  mem_request = ""
  group_resource_type = ""
  subnet_id = ""
  agent_cpu_request = ""
  agent_cpu_limit = ""
  agent_mem_request = ""
  agent_mem_limit = ""
  istio_cpu_request = ""
  istio_cpu_limit = ""
  istio_mem_request = ""
  istio_mem_limit = ""
                                          }

`
