package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwVpcFwSwitchesDataSource_basic -v
func TestAccTencentCloudNeedFixCfwVpcFwSwitchesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwVpcFwSwitchesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cfw_vpc_fw_switches.example"),
				),
			},
		},
	})
}

const testAccCfwVpcFwSwitchesDataSource = `
data "tencentcloud_cfw_vpc_fw_switches" "example" {
  vpc_ins_id = "cfwg-c8c2de41"
}
`
