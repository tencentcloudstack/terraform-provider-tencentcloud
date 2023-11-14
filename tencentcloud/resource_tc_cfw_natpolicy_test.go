package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCfwNatpolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwNatpolicy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfw_natpolicy.natpolicy", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfw_natpolicy.natpolicy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfwNatpolicy = `

resource "tencentcloud_cfw_natpolicy" "natpolicy" {
  rules {
		source_content = "192.168.0.2"
		source_type = "ip"
		target_content = "192.168.0.2"
		target_type = "ip"
		protocol = "TCP"
		rule_action = "allow"
		port = "80"
		direction = 1
		order_index = 1
		enable = "true"
		uuid = 1
		description = "test"

  }
  from = ""
}

`
