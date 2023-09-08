package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWafTlsVersionsDataSource_basic -v
func TestAccTencentCloudNeedFixWafTlsVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafTlsVersionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_tls_versions.example"),
				),
			},
		},
	})
}

const testAccWafTlsVersionsDataSource = `
data "tencentcloud_waf_tls_versions" "example" {}
`
