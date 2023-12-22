package gaap_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcgaap "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/gaap"

	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_gaap_security_rule
	resource.AddTestSweepers("tencentcloud_gaap_security_rule", &resource.Sweeper{
		Name: "tencentcloud_gaap_security_rule",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			sharedClient, err := tcacctest.SharedClientForRegion(r)
			if err != nil {
				return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
			}
			client := sharedClient.(tccommon.ProviderMeta)
			service := svcgaap.NewGaapService(client.GetAPIV3Conn())

			securityRules, err := service.DescribeSecurityRules(ctx, tcacctest.DefaultGaapSecurityPolicyId)
			if err != nil {
				return err
			}
			for _, securityRule := range securityRules {
				instanceName := *securityRule.AliasName

				if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
					continue
				}

				ee := service.DeleteSecurityRule(ctx, tcacctest.DefaultGaapSecurityPolicyId, *securityRule.RuleId)
				if ee != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudGaapSecurityRuleResource_basic(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_security_rule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRuleResource_drop(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleDrop,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "DROP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRuleResource_name(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleWithName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "81"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "name", "ci-test-gaap-sr"),
				),
			},
			{
				Config: testAccGaapSecurityRuleUpdateName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "name", "ci-test-gaap-sr-new"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRuleResource_ipSubnet(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleIpSubnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "80"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRuleResource_allProtocols(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleAllProtocols,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "ALL"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityRuleResource_AllPorts(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapSecurityRuleDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityRuleAllPorts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityRuleExists("tencentcloud_gaap_security_rule.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_rule.foo", "port", "ALL"),
				),
			},
		},
	})
}

func testAccCheckGaapSecurityRuleExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no security rule ID is set")
		}

		service := svcgaap.NewGaapService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rule, err := service.DescribeSecurityRule(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if rule == nil {
			return errors.New("no security rule not exist")
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckGaapSecurityRuleDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		service := svcgaap.NewGaapService(client)

		rule, err := service.DescribeSecurityRule(context.TODO(), *id)
		if err != nil {
			return err
		}

		if rule != nil {
			return errors.New("security rule still exists")
		}

		return nil
	}
}

var testAccGaapSecurityRuleBasic = fmt.Sprintf(`
resource tencentcloud_gaap_security_rule "foo" {
  policy_id = "%s"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}
`, tcacctest.DefaultGaapSecurityPolicyId)

var testAccGaapSecurityRuleWithName = fmt.Sprintf(`
resource tencentcloud_gaap_security_rule "foo" {
  name      = "ci-test-gaap-sr"
  policy_id = "%s"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "81"
}
`, tcacctest.DefaultGaapSecurityPolicyId)

var testAccGaapSecurityRuleDrop = fmt.Sprintf(`

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = "%s"
  cidr_ip   = "1.1.1.1"
  action    = "DROP"
  protocol  = "TCP"
  port      = "80"
}
`, tcacctest.DefaultGaapSecurityPolicyId)

var testAccGaapSecurityRuleUpdateName = fmt.Sprintf(`
resource tencentcloud_gaap_security_rule "foo" {
  name      = "ci-test-gaap-sr-new"
  policy_id = "%s"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "81"
}
`, tcacctest.DefaultGaapSecurityPolicyId)

var testAccGaapSecurityRuleIpSubnet = fmt.Sprintf(`

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = "%s"
  cidr_ip   = "192.168.1.0/24"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}
`, tcacctest.DefaultGaapSecurityPolicyId)

var testAccGaapSecurityRuleAllProtocols = fmt.Sprintf(`

resource tencentcloud_gaap_security_rule "foo" {
  policy_id = "%s"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
}
`, tcacctest.DefaultGaapSecurityPolicyId)

var testAccGaapSecurityRuleAllPorts = fmt.Sprintf(`
resource tencentcloud_gaap_security_rule "foo" {
  policy_id = "%s"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
}
`, tcacctest.DefaultGaapSecurityPolicyId)
