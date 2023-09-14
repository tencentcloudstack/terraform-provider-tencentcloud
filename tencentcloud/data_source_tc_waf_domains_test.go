package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWafDomainsDataSource_basic -v
func TestAccTencentCloudNeedFixWafDomainsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomainsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_domains.example"),
				),
			},
		},
	})
}

const testAccWafDomainsDataSource = `
data "tencentcloud_waf_domains" "example" {
  instance_id = "waf_2kxtlbky01b3wceb"
  domain      = "tf.example.com"
}
`
