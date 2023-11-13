package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWafFindDomainsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafFindDomainsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_find_domains.find_domains")),
			},
		},
	})
}

const testAccWafFindDomainsDataSource = `

data "tencentcloud_waf_find_domains" "find_domains" {
  key = ""
  is_waf_domain = ""
  by = ""
  order = ""
  }

`
