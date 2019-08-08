package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudGaapSecurityRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRule_withName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleWithName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "name", "ci-test-gaap-sr"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRule_drop(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleDrop,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "DROP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRule_updateName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleWithName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "name", "ci-test-gaap-sr"),
				),
			},
			{
				Config: testAccGaapSecurityRuleUpdateName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "name", "ci-test-gaap-sr-new"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRule_ipSubnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleIpSubnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRule_allProtocols(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleAllProtocols,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRule_AllPorts(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleAllPorts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_security_rule.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "ALL"),
				),
			},
		},
	})
}

const testAccGaapSecurityRuleBasic = `
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

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}
`

const testAccGaapSecurityRuleWithName = `
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

resource tencentcloud_gaap_security_rule "foo" {
  name      = "ci-test-gaap-sr"
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}
`

const testAccGaapSecurityRuleDrop = `
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

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  cidr_ip   = "1.1.1.1"
  action    = "DROP"
  protocol  = "TCP"
  port      = "80"
}
`

const testAccGaapSecurityRuleUpdateName = `
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

resource tencentcloud_gaap_security_rule "foo" {
  name      = "ci-test-gaap-sr-new"
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}
`

const testAccGaapSecurityRuleIpSubnet = `
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

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  cidr_ip   = "192.168.1.0/24"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}
`

const testAccGaapSecurityRuleAllProtocols = `
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

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  port      = "80"
}
`

const testAccGaapSecurityRuleAllPorts = `
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

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
}
`
