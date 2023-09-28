package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwNatFwSwitchesDataSource_basic -v
func TestAccTencentCloudNeedFixCfwNatFwSwitchesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwNatFwSwitchesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cfw_nat_fw_switches.example"),
				),
			},
		},
	})
}

const testAccCfwNatFwSwitchesDataSource = `
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
}
`
