package cfw_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwEdgeFwSwitchesDataSource_basic -v
func TestAccTencentCloudNeedFixCfwEdgeFwSwitchesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwEdgeFwSwitchesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cfw_edge_fw_switches.example"),
				),
			},
		},
	})
}

const testAccCfwEdgeFwSwitchesDataSource = `
data "tencentcloud_cfw_edge_fw_switches" "example" {}
`
