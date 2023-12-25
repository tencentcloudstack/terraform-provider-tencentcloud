package vpc_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudSecurityGroupRule_basic(t *testing.T) {
	t.Parallel()
	var sgrId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy(&sgrId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.http-in", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.http-in", "cidr_ip", "1.1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.http-in", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.http-in", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.http-in", "type", "ingress"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.http-in", "policy_index", "0"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.http-in", "source_sgid"),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupRule_ssh(t *testing.T) {
	t.Parallel()
	var sgrId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy(&sgrId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfigSSH,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.ssh-in", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.ssh-in", "cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.ssh-in", "ip_protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.ssh-in", "port_range", "22"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.ssh-in", "description", "ssh in rule"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.ssh-in", "source_sgid"),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupRule_egress(t *testing.T) {
	t.Parallel()
	var sgrId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy(&sgrId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfigEgress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.egress-drop", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "cidr_ip", "10.2.3.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "ip_protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "port_range", "3000-4000"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "policy", "DROP"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "type", "EGRESS"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.egress-drop", "source_sgid"),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupRule_sourcesgid(t *testing.T) {
	t.Parallel()
	var sgrId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy(&sgrId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfigSourceSGID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.sourcesgid-in", &sgrId),
					resource.TestCheckResourceAttrSet("tencentcloud_security_group_rule.sourcesgid-in", "source_sgid"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.sourcesgid-in", "ip_protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.sourcesgid-in", "port_range", "80,8080"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.sourcesgid-in", "policy", "ACCEPT"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.sourcesgid-in", "cidr_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupRule_allDrop(t *testing.T) {
	t.Parallel()
	var sgrId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy(&sgrId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfigAllDrop,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.egress-drop", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "ip_protocol", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "port_range", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "policy", "DROP"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.egress-drop", "source_sgid"),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupRule_addressTemplate(t *testing.T) {
	t.Parallel()
	var sgrId string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy(&sgrId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfigAddressTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.address_template", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.address_template", "address_template.#", "1"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.address_template", "address_template.group_id"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.address_template", "port_range", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.address_template", "ip_protocol", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.address_template", "policy", "DROP"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.address_template", "cidr_ip"),
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.address_template_group", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.address_template_group", "address_template.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.address_template_group", "ip_protocol", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.address_template_group", "policy", "DROP"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.address_template_group", "cidr_ip"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.address_template_group", "port_range", "ALL"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.address_template_group", "address_template.template_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupRule_protocolTemplate(t *testing.T) {
	t.Parallel()
	var sgrId string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy(&sgrId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfigProtocolTemplate_multi_rule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.protocol_template1", &sgrId),
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.protocol_template2", &sgrId),
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.protocol_template3", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.protocol_template1", "protocol_template.#", "1"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.protocol_template1", "protocol_template.group_id"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.protocol_template1", "cidr_ip", "0.0.0.0/0"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.protocol_template1", "ip_protocol"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.protocol_template1", "policy", "DROP"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.protocol_template2", "policy", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.protocol_template3", "policy", "ACCEPT"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.protocol_template1", "cip_protocol"),
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.protocol_template_group", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.protocol_template_group", "protocol_template.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.protocol_template_group", "cidr_ip", "0.0.0.0/0"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.protocol_template_group", "port_range"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.protocol_template_group", "policy", "DROP"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.protocol_template_group", "ip_protocol"),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group_rule.protocol_template_group", "protocol_template.template_id"),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupRuleDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, _, policy, err := service.DescribeSecurityGroupPolicy(context.TODO(), *id)
		if err != nil {
			return err
		}

		if policy == nil {
			return nil
		}

		return errors.New("security group rule still exist")
	}
}

func testAccCheckSecurityGroupRuleExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no security group rule ID is set")
		}

		*id = rs.Primary.ID

		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, _, policy, err := service.DescribeSecurityGroupPolicy(context.TODO(), *id)
		if err != nil {
			return err
		}

		if policy == nil {
			return errors.New("security group rule not exist")
		}

		return nil
	}
}

const testAccSecurityGroupRuleConfig = `
resource "tencentcloud_security_group" "foo" {
  name        = "ci-temp-test-sg"
  description = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_rule" "http-in" {
  security_group_id = tencentcloud_security_group.foo.id
  type              = "ingress"
  cidr_ip           = "1.1.1.1"
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
  policy_index      = 0
}
`

const testAccSecurityGroupRuleConfigSSH = `
resource "tencentcloud_security_group" "foo" {
  name        = "ci-temp-test-sg"
  description = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_rule" "ssh-in" {
  security_group_id = tencentcloud_security_group.foo.id
  type              = "INGRESS"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "TCP"
  port_range        = "22"
  policy            = "ACCEPT"
  description       = "ssh in rule"
}
`

const testAccSecurityGroupRuleConfigEgress = `
resource "tencentcloud_security_group" "foo" {
  name        = "ci-temp-test-sg"
  description = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_rule" "egress-drop" {
  security_group_id = tencentcloud_security_group.foo.id
  type              = "EGRESS"
  cidr_ip           = "10.2.3.0/24"
  ip_protocol       = "UDP"
  port_range        = "3000-4000"
  policy            = "DROP"
}
`

const testAccSecurityGroupRuleConfigSourceSGID = `
resource "tencentcloud_security_group" "foo" {
  name        = "ci-temp-test-sg"
  description = "ci-temp-test-sg"
}

resource "tencentcloud_security_group" "boo" {
  name        = "ci-temp-test-sg"
  description = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_rule" "sourcesgid-in" {
  security_group_id = tencentcloud_security_group.foo.id
  type              = "ingress"
  source_sgid		= tencentcloud_security_group.boo.id
  ip_protocol       = "TCP"
  port_range        = "80,8080"
  policy            = "ACCEPT"
}
`

const testAccSecurityGroupRuleConfigAllDrop = `
resource "tencentcloud_security_group" "foo" {
  name        = "ci-temp-test-sg"
  description = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_rule" "egress-drop" {
  security_group_id = tencentcloud_security_group.foo.id
  cidr_ip           = "0.0.0.0/0"
  type              = "ingress"
  policy            = "DROP"
}
`

const testAccSecurityGroupRuleConfigAddressTemplate = `
resource "tencentcloud_security_group" "foo" {
  name        = "ci-temp-test-sg"
  description = "ci-temp-test-sg"
}

resource "tencentcloud_address_template" "templateB" {
  name = "testB"
  addresses = ["1.1.1.1/24", "1.1.1.0-1.1.1.1"]
}

resource "tencentcloud_address_template_group" "group"{
	name = "test_update"
	template_ids = [tencentcloud_address_template.templateB.id]
}

resource "tencentcloud_security_group_rule" "address_template_group" {
  security_group_id = tencentcloud_security_group.foo.id
  type              = "ingress"
  policy            = "DROP"

  address_template  {
		group_id = tencentcloud_address_template_group.group.id
	}
}

resource "tencentcloud_security_group_rule" "address_template" {
  security_group_id = tencentcloud_security_group.foo.id
  type              = "INGRESS"
  policy            = "DROP"

  address_template  {
		template_id = tencentcloud_address_template.templateB.id
	}
}
`

const testAccSecurityGroupRuleConfigProtocolTemplate_multi_rule = `
resource "tencentcloud_security_group" "foo" {
  name        = "ci-temp-test-sg"
  description = "ci-temp-test-sg"
}

resource "tencentcloud_protocol_template" "templateB" {
  name = "testB"
  protocols = ["tcp:80", "udp:90,111"]
}

resource "tencentcloud_protocol_template_group" "group"{
	name = "test_update"
	template_ids = [tencentcloud_protocol_template.templateB.id]
}

resource "tencentcloud_security_group_rule" "protocol_template_group" {
  security_group_id = tencentcloud_security_group.foo.id
  type              = "ingress"
  policy            = "DROP"
  cidr_ip           = "0.0.0.0/0"

  protocol_template  {
		group_id = tencentcloud_protocol_template_group.group.id
	}
}

resource "tencentcloud_security_group_rule" "protocol_template1" {
  security_group_id = tencentcloud_security_group.foo.id
  type              = "INGRESS"
  policy            = "DROP"
  cidr_ip           = "0.0.0.0/0"

  protocol_template  {
		template_id = tencentcloud_protocol_template.templateB.id
	}
}

resource "tencentcloud_security_group_rule" "protocol_template2" {
	security_group_id = tencentcloud_security_group.foo.id
	type              = "INGRESS"
	policy            = "ACCEPT"
	cidr_ip           = "0.0.0.0/0"
  
	protocol_template  {
		  template_id = tencentcloud_protocol_template.templateB.id
	  }
  }

  resource "tencentcloud_security_group_rule" "protocol_template3" {
	security_group_id = tencentcloud_security_group.foo.id
	type              = "INGRESS"
	policy            = "ACCEPT"
	cidr_ip           = "10.0.0.0/12"
  
	protocol_template  {
		  template_id = tencentcloud_protocol_template.templateB.id
	  }
  }
`
