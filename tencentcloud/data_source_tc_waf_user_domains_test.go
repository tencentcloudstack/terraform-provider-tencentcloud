package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWafUserDomainsDataSource_basic -v
func TestAccTencentCloudNeedFixWafUserDomainsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafUserDomainsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_user_domains.example"),
				),
			},
		},
	})
}

const testAccWafUserDomainsDataSource = `
data "tencentcloud_waf_user_domains" "example" {}
`
