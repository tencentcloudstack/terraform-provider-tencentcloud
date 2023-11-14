package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbSwitchHAResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbSwitchHA,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_switch_h_a.switch_h_a", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_switch_h_a.switch_h_a",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbSwitchHA = `

resource "tencentcloud_mariadb_switch_h_a" "switch_h_a" {
  instance_id = ""
  zone = ""
}

`
