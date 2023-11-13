package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbOpenDBExtranetAccessResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbOpenDBExtranetAccess,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_open_d_b_extranet_access.open_d_b_extranet_access", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_open_d_b_extranet_access.open_d_b_extranet_access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbOpenDBExtranetAccess = `

resource "tencentcloud_dcdb_open_d_b_extranet_access" "open_d_b_extranet_access" {
  instance_id = ""
  ipv6_flag = 
}

`
