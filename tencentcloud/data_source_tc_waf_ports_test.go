package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWafPortsDataSource_basic -v
func TestAccTencentCloudNeedFixWafPortsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPortsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_ports.example"),
				),
			},
		},
	})
}

const testAccWafPortsDataSource = `
data "tencentcloud_waf_ports" "example" {
  edition     = "clb-waf"
  instance_id = "waf_2kxtlbky00b2v1fn"
}
`
