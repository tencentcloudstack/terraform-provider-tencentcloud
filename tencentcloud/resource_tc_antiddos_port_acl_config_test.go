package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosPortAclConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosPortAclConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_port_acl_config.port_acl_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_port_acl_config.port_acl_config", "instance_id", "bgp-00000ry7"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_port_acl_config.port_acl_config", "acl_config.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_antiddos_port_acl_config.port_acl_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosPortAclConfig = `
resource "tencentcloud_antiddos_port_acl_config" "port_acl_config" {
	instance_id = "bgp-00000ry7"
	acl_config {
	  forward_protocol = "all"
	  d_port_start     = 22
	  d_port_end       = 23
	  s_port_start     = 22
	  s_port_end       = 23
	  action           = "drop"
	  priority         = 2
  
	}
  }
`
