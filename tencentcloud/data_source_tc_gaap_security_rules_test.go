package tencentcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_security_rules.foo", "rules.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.foo", "rules.0.action"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapSecurityRules_multi(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesRuleId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.ruleId"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.ruleId", "rules.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.ruleId", "rules.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.ruleId", "rules.0.cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.ruleId", "rules.0.name", "ci-test-gaap-s-rule"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.ruleId", "rules.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.ruleId", "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.ruleId", "rules.0.action", "ACCEPT"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesAction,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.action"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_security_rules.action", "rules.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.action", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.action", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.action", "rules.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.action", "rules.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.action", "rules.0.protocol"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.action", "rules.0.action", "ACCEPT"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesCidrIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.cidrIp"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_security_rules.cidrIp", "rules.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.cidrIp", "rules.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.cidrIp", "rules.0.cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.cidrIp", "rules.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.cidrIp", "rules.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.cidrIp", "rules.0.protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.cidrIp", "rules.0.action"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.name"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_security_rules.name", "rules.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.name", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.name", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.name", "rules.0.name", "ci-test-gaap-s-rule"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.name", "rules.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.name", "rules.0.protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.name", "rules.0.action"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.port"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_security_rules.port", "rules.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.port", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.port", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.port", "rules.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.port", "rules.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.port", "rules.0.protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.port", "rules.0.action"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityRulesProtocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_rules.protocol"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_security_rules.protocol", "rules.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.protocol", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.protocol", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.protocol", "rules.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.protocol", "rules.0.port"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_rules.protocol", "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_rules.protocol", "rules.0.action"),
				),
			},
		},
	})
}

var TestAccDataSourceTencentCloudGaapSecurityRulesBasic = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "foo" {
  policy_id = tencentcloud_gaap_security_rule.foo.policy_id
}
`, defaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapSecurityRulesRuleId = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "ruleId" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  rule_id   = tencentcloud_gaap_security_rule.foo.id
}
`, defaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapSecurityRulesAction = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "action" {
  policy_id = tencentcloud_gaap_security_rule.foo.policy_id
  action    = "ACCEPT"
}
`, defaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapSecurityRulesCidrIp = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "cidrIp" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  cidr_ip   = tencentcloud_gaap_security_rule.foo.cidr_ip
}
`, defaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapSecurityRulesName = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "name" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  name      = tencentcloud_gaap_security_rule.foo.name
}
`, defaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapSecurityRulesPort = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "port" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  port      = tencentcloud_gaap_security_rule.foo.port
}
`, defaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapSecurityRulesProtocol = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data tencentcloud_gaap_security_rules "protocol" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  protocol  = tencentcloud_gaap_security_rule.foo.protocol
}
`, defaultGaapProxyId)
