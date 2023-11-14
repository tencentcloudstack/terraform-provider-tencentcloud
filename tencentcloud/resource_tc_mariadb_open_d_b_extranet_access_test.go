package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbOpenDBExtranetAccessResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbOpenDBExtranetAccess,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_open_d_b_extranet_access.open_d_b_extranet_access", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_open_d_b_extranet_access.open_d_b_extranet_access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbOpenDBExtranetAccess = `

resource "tencentcloud_mariadb_open_d_b_extranet_access" "open_d_b_extranet_access" {
  instance_id = ""
  ipv6_flag = 
}

`
