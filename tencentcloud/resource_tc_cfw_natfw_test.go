package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCfwNatfwResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwNatfw,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfw_natfw.natfw", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfw_natfw.natfw",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfwNatfw = `

resource "tencentcloud_cfw_natfw" "natfw" {
  name = "test natfw"
  width = 20
  mode = 0
  new_mode_items {
		vpc_list = 
		eips = 
		add_count = 1

  }
  nat_gw_list = 
  zone = "ap-guangzhou-1"
  zone_bak = "ap-guangzhou-2"
  cross_a_zone = 1
  is_create_domain = 0
  domain = ""
  fw_cidr_info {
		fw_cidr_type = "VpcSelf"
		fw_cidr_lst {
			vpc_id = "vpc-id"
			fw_cidr = "10.96.0.0/16"
		}
		com_fw_cidr = ""

  }
}

`
