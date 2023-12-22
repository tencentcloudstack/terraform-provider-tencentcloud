package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWafPortsDataSource_basic -v
func TestAccTencentCloudNeedFixWafPortsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPortsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_waf_ports.example"),
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
