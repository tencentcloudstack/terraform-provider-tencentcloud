package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafUserClbRegionsDataSource_basic -v
func TestAccTencentCloudWafUserClbRegionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafUserClbRegionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_waf_user_clb_regions.example"),
				),
			},
		},
	})
}

const testAccWafUserClbRegionsDataSource = `
data "tencentcloud_waf_user_clb_regions" "example" {}
`
