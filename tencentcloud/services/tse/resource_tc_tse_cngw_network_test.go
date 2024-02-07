package tse_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudNeedFixTseCngwNetworkResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwNetwork,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_network.cngw_network", "id")),
			},
		},
	})
}

const testAccTseCngwNetwork = `

resource "tencentcloud_tse_cngw_network" "cngw_network" {
  gateway_id                 = "gateway-cf8c99c3"
  group_id                   = "group-a160d123"
  internet_address_version   = "IPV4"
  internet_pay_mode          = "BANDWIDTH"
  description                = "des-test1"
  internet_max_bandwidth_out = 1
  master_zone_id             = "ap-guangzhou-3"
  multi_zone_flag            = true
  sla_type                   = "clb.c2.medium"
  slave_zone_id              = "ap-guangzhou-4"
}

`
