package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbCloseDBExtranetAccessResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbCloseDBExtranetAccess,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_close_d_b_extranet_access.close_d_b_extranet_access", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_close_d_b_extranet_access.close_d_b_extranet_access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbCloseDBExtranetAccess = `

resource "tencentcloud_dcdb_close_d_b_extranet_access" "close_d_b_extranet_access" {
  instance_id = ""
  ipv6_flag = 
}

`
