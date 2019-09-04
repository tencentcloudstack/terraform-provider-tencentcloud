package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudGaapSecurityPolices_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityPolicesBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_policies.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_policies.foo", "proxy_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_policies.foo", "status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_policies.foo", "action", "ACCEPT"),
				),
			},
		},
	})
}

var TestAccDataSourceTencentCloudGaapSecurityPolicesBasic = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

data tencentcloud_gaap_security_policies "foo" {
  id = "${tencentcloud_gaap_security_policy.foo.id}"
}
`, GAAP_PROXY_ID)
