package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwEdgeFwSwitchesDataSource_basic -v
func TestAccTencentCloudNeedFixCfwEdgeFwSwitchesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwEdgeFwSwitchesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cfw_edge_fw_switches.example"),
				),
			},
		},
	})
}

const testAccCfwEdgeFwSwitchesDataSource = `
data "tencentcloud_cfw_edge_fw_switches" "example" {}
`
