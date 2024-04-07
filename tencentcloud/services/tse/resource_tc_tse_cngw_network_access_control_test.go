package tse_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudNeedFixTseCngwNetworkAccessControlResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwNetworkAccessControl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_network_access_control.cngw_network_access_control", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_network_access_control.cngw_network_access_control", "access_control.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_network_access_control.cngw_network_access_control", "access_control.0.mode", "Whitelist"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_network_access_control.cngw_network_access_control", "access_control.0.cidr_white_list.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_network_access_control.cngw_network_access_control",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTseCngwNetworkAccessControl = `

resource "tencentcloud_tse_cngw_network_access_control" "cngw_network_access_control" {
	gateway_id                 = "gateway-cf1790c7"
	group_id                   = "group-d8d99615"
	network_id                 = "network-9cd9821f"
	access_control {
	  mode            = "Whitelist"
	  cidr_white_list = ["1.1.1.0"]
	}
  }

`
