package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudGaapSecurityPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_policy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_security_policy.foo", "status"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityPolicy_disable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_policy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_security_policy.foo", "status"),
				),
			},
			{
				Config: testAccGaapSecurityPolicyDisable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_policy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_security_policy.foo", "status"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityPolicy_drop(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityPolicyDrop,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_policy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "action", "DROP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_security_policy.foo", "status"),
				),
			},
		},
	})
}

const testAccGaapSecurityPolicyBasic = `
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
`

const testAccGaapSecurityPolicyDisable = `
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
  enable   = false
}
`

const testAccGaapSecurityPolicyDrop = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  action   = "DROP"
}
`
