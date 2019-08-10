package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudGaapSecurityRules_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.rule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.name", "ci-test-gaap-s-rule"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.action", "ACCEPT"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapSecurityRules_ruleId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesRuleId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.rule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.name", "ci-test-gaap-s-rule"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.action", "ACCEPT"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapSecurityRules_action(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesAction,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.rule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.name", "ci-test-gaap-s-rule"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.action", "ACCEPT"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapSecurityRules_cidrIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesCidrIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.rule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.name", "ci-test-gaap-s-rule"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.action", "ACCEPT"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapSecurityRules_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.rule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.name", "ci-test-gaap-s-rule"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.action", "ACCEPT"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapSecurityRules_port(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.rule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.name", "ci-test-gaap-s-rule"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.action", "ACCEPT"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapSecurityRules_protocol(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesProtocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.rule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.name", "ci-test-gaap-s-rule"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.0.action", "ACCEPT"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudGaapSecurityRulesBasic = `
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
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
}
`

const TestAccDataSourceTencentCloudGaapSecurityRulesRuleId = `
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
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

resource tencentcloud_gaap_security_rule "bar" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "UDP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  rule_id   = "${tencentcloud_gaap_security_rule.foo.id}"
}
`

const TestAccDataSourceTencentCloudGaapSecurityRulesAction = `
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
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

resource tencentcloud_gaap_security_rule "bar" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "DROP"
  protocol  = "UDP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  action    = "ACCEPT"
}
`

const TestAccDataSourceTencentCloudGaapSecurityRulesCidrIp = `
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
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

resource tencentcloud_gaap_security_rule "bar" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "192.168.1.0/24"
  action    = "DROP"
  protocol  = "UDP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  cidr_ip   = "1.1.1.1"
}
`

const TestAccDataSourceTencentCloudGaapSecurityRulesName = `
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
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

resource tencentcloud_gaap_security_rule "bar" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  name      = "ci-test-gaap-s-rule-2"
  cidr_ip   = "192.168.1.0/24"
  action    = "DROP"
  protocol  = "UDP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  name      = "ci-test-gaap-s-rule"
}
`

const TestAccDataSourceTencentCloudGaapSecurityRulesPort = `
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
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

resource tencentcloud_gaap_security_rule "bar" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  name      = "ci-test-gaap-s-rule-2"
  cidr_ip   = "192.168.1.0/24"
  action    = "DROP"
  protocol  = "UDP"
  port      = "8080"
}

data tencentcloud_gaap_security_rules "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  port      = "80"
}
`

const TestAccDataSourceTencentCloudGaapSecurityRulesProtocol = `
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
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

resource tencentcloud_gaap_security_rule "bar" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  name      = "ci-test-gaap-s-rule-2"
  cidr_ip   = "192.168.1.0/24"
  action    = "DROP"
  protocol  = "UDP"
  port      = "8080"
}

data tencentcloud_gaap_security_rules "foo" {
  policy_id = "${tencentcloud_gaap_security_policy.foo.id}"
  protocol  = "TCP"
}
`
