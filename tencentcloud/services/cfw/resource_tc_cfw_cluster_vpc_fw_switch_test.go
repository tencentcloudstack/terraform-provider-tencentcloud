package cfw_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCfwClusterVpcFwSwitchResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwClusterVpcFwSwitch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_cluster_vpc_fw_switch.example", "id"),
				),
			},
			{
				Config: testAccCfwClusterVpcFwSwitchUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_cluster_vpc_fw_switch.example", "id"),
				),
			},
		},
	})
}

const testAccCfwClusterVpcFwSwitch = `
resource "tencentcloud_cfw_cluster_vpc_fw_switch" "example" {
  ccn_id       = "ccn-8qv0ro89"
  switch_mode  = 2
  routing_mode = 0
  region_cidr_configs {
    region      = "ap-guangzhou"
    cidr_mode   = 1
    custom_cidr = ""
  }
}
`

const testAccCfwClusterVpcFwSwitchUpdate = `
resource "tencentcloud_cfw_cluster_vpc_fw_switch" "example" {
  ccn_id       = "ccn-8qv0ro89"
  switch_mode  = 2
  routing_mode = 0
  region_cidr_configs {
    region      = "ap-chongqing"
    cidr_mode   = 1
    custom_cidr = ""
  }
}
`
