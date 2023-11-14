package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCfwNatFwSwitchDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwNatFwSwitchDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cfw_nat_fw_switch.nat_fw_switch")),
			},
		},
	})
}

const testAccCfwNatFwSwitchDataSource = `

data "tencentcloud_cfw_nat_fw_switch" "nat_fw_switch" {
  search_value = ""
  status = 
  vpc_id = ""
  nat_id = ""
  nat_ins_id = ""
  area = ""
  }

`
