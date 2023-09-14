package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWafFindDomainsDataSource_basic -v
func TestAccTencentCloudNeedFixWafFindDomainsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafFindDomainsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_find_domains.example"),
				),
			},
		},
	})
}

const testAccWafFindDomainsDataSource = `
data "tencentcloud_waf_find_domains" "example" {
  key           = "keyWord"
  is_waf_domain = "1"
  by            = "FindTime"
  order         = "asc"
}
`
