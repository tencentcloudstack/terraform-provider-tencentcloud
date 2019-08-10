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
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_polices.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_polices.foo", "proxy_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_polices.foo", "status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_polices.foo", "action", "ACCEPT"),
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
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  action   = "ACCEPT"
}

data tencentcloud_gaap_security_polices "foo" {
  id = "${tencentcloud_gaap_security_policy.foo.id}"
}
`
