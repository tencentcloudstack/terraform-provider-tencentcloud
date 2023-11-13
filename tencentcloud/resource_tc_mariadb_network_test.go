package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbNetworkResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbNetwork,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_network.network", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_network.network",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbNetwork = `

resource "tencentcloud_mariadb_network" "network" {
  instance_id = ""
  vpc_id = ""
  subnet_id = ""
  vip = ""
  vipv6 = ""
  vip_release_delay = 
}

`
