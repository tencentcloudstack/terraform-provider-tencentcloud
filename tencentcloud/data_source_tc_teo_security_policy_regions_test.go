package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoSecurityPolicyRegionsDataSource -v
func TestAccTencentCloudTeoSecurityPolicyRegionsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeoSecurityPolicyRegions,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_teo_security_policy_regions.security_policy_regions"),
				),
			},
		},
	})
}

const testAccDataSourceTeoSecurityPolicyRegions = `

data "tencentcloud_teo_security_policy_regions" "security_policy_regions" {
}

`
