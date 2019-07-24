package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudSecurityGroupRule_basic(t *testing.T) {
	var sgrId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy(&sgrId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.http-in", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.http-in", "cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.http-in", "ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.http-in", "description", ""),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupRule_ssh(t *testing.T) {
	var sgrId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy(&sgrId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfigSSH,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.ssh-in", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.ssh-in", "cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.ssh-in", "ip_protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.ssh-in", "port_range", "22"),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupRule_egress(t *testing.T) {
	var sgrId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy(&sgrId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfigEgress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.egress-DROP", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-DROP", "cidr_ip", "10.2.3.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-DROP", "ip_protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-DROP", "port_range", "3000-4000"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-DROP", "policy", "DROP"),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupRule_sourcesgid(t *testing.T) {
	var sgrId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupRule_allDrop(t *testing.T) {
	var sgrId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupRuleDestroy(&sgrId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleConfigAllDrop,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.egress-DROP", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-DROP", "cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-DROP", "ip_protocol", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-DROP", "port_range", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-DROP", "policy", "DROP"),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupRuleDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

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

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

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
  security_group_id = "${tencentcloud_security_group.foo.id}"
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
}
`

const testAccSecurityGroupRuleConfigSSH = `
resource "tencentcloud_security_group" "foo" {
  name        = "ci-temp-test-sg"
  description = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_rule" "ssh-in" {
  security_group_id = "${tencentcloud_security_group.foo.id}"
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

resource "tencentcloud_security_group_rule" "egress-DROP" {
  security_group_id = "${tencentcloud_security_group.foo.id}"
  type              = "egress"
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
  security_group_id = "${tencentcloud_security_group.foo.id}"
  type              = "ingress"
  source_sgid		= "${tencentcloud_security_group.boo.id}"
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

resource "tencentcloud_security_group_rule" "egress-DROP" {
  security_group_id = "${tencentcloud_security_group.foo.id}"
  cidr_ip           = "0.0.0.0/0"
  type              = "ingress"
  policy            = "DROP"
}
`
