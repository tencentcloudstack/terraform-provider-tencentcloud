package tencentcloud

import (
	"encoding/json"
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
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.ssh-in", "ip_protocol", "tcp"),
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
					testAccCheckSecurityGroupRuleExists("tencentcloud_security_group_rule.egress-drop", &sgrId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "cidr_ip", "10.2.3.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "ip_protocol", "udp"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "port_range", "3000-4000"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.egress-drop", "policy", "drop"),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupRuleDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*TencentCloudClient).commonConn

		rule, _ := parseSecurityGroupRuleId(*id)

		params := map[string]string{
			"Action":    "DescribeSecurityGroupEx",
			"projectId": "0",
			"sgId":      rule["sgId"],
		}
		response, err := conn.SendRequest("dfw", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Code     int    `json:"code"`
			CodeDesc string `json:"codeDesc"`
			Data     struct {
				TotalNum int `json:"totalNum"`
			} `json:"data"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return err
		}
		if jsonresp.Data.TotalNum == 0 {
			return nil
		}

		_, err = describeSecurityGroupRuleIndex(conn, rule)
		if err == errSecurityGroupRuleNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		return fmt.Errorf("Security group rule still exists.")
	}
}

func testAccCheckSecurityGroupRuleExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No security group ID is set")
		}

		conn := testAccProvider.Meta().(*TencentCloudClient).commonConn
		rule, _ := parseSecurityGroupRuleId(rs.Primary.ID)
		_, err := describeSecurityGroupRuleIndex(conn, rule)
		if err != nil {
			return err
		}

		*id = rs.Primary.ID
		return nil
	}
}

const testAccSecurityGroupRuleConfig = `
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
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
  name = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_rule" "ssh-in" {
  security_group_id = "${tencentcloud_security_group.foo.id}"
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "22"
  policy            = "accept"
}
`

const testAccSecurityGroupRuleConfigEgress = `
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_rule" "egress-drop" {
  security_group_id = "${tencentcloud_security_group.foo.id}"
  type              = "egress"
  cidr_ip           = "10.2.3.0/24"
  ip_protocol       = "udp"
  port_range        = "3000-4000"
  policy            = "drop"
}
`
