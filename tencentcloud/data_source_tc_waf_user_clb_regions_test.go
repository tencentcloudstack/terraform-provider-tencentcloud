package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafUserClbRegionsDataSource_basic -v
func TestAccTencentCloudWafUserClbRegionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafUserClbRegionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_user_clb_regions.example"),
				),
			},
		},
	})
}

const testAccWafUserClbRegionsDataSource = `
data "tencentcloud_waf_user_clb_regions" "example" {}
`
