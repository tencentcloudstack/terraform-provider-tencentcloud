package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbCloseDBExtranetAccessResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbCloseDBExtranetAccess,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_close_d_b_extranet_access.close_d_b_extranet_access", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_close_d_b_extranet_access.close_d_b_extranet_access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbCloseDBExtranetAccess = `

resource "tencentcloud_mariadb_close_d_b_extranet_access" "close_d_b_extranet_access" {
  instance_id = ""
  ipv6_flag = 
}

`
