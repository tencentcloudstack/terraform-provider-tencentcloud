package tencentcloud

import (
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

const TestAccDataSourceTencentCloudGaapSecurityPolicesBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  action   = "ACCEPT"
}

data tencentcloud_gaap_security_policies "foo" {
  id = "${tencentcloud_gaap_security_policy.foo.id}"
}
`
