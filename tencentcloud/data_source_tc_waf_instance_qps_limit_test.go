package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafInstanceQpsLimitDataSource_basic -v
func TestAccTencentCloudWafInstanceQpsLimitDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafInstanceQpsLimitDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_instance_qps_limit.example"),
				),
			},
		},
	})
}

const testAccWafInstanceQpsLimitDataSource = `
data "tencentcloud_waf_instance_qps_limit" "example" {
  instance_id = "waf_2kxtlbky00b3b4qz"
}
`
