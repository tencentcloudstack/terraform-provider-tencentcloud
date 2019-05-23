package tencentcloud

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudSecurityGroup_basic(t *testing.T) {
	var sgId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy(&sgId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("tencentcloud_security_group.foo", &sgId),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "name", "ci-temp-test-sg"),
				),
			},
			{
				ResourceName:      "tencentcloud_security_group.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudSecurityGroup_update(t *testing.T) {
	var sgId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy(&sgId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("tencentcloud_security_group.foo", &sgId),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "name", "ci-temp-test-sg"),
				),
			},
			{
				Config: testAccSecurityGroupConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("tencentcloud_security_group.foo", &sgId),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "name", "ci-temp-test-sg-updated"),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*TencentCloudClient).commonConn
		params := map[string]string{
			"Action":    "DescribeSecurityGroupEx",
			"projectId": "0",
			"sgId":      *id,
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
		if jsonresp.Data.TotalNum != 0 {
			return fmt.Errorf("Security group still exists.")
		}
		return nil
	}
}

func testAccCheckSecurityGroupExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No security group ID is set")
		}

		client := testAccProvider.Meta().(*TencentCloudClient).commonConn
		params := map[string]string{
			"Action":    "DescribeSecurityGroupEx",
			"projectId": "0",
			"sgId":      rs.Primary.ID,
		}
		response, err := client.SendRequest("dfw", params)
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
			return fmt.Errorf("Security group not found: %s", rs.Primary.ID)
		}

		*id = rs.Primary.ID
		return nil
	}
}

const testAccSecurityGroupConfig = `
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
}
`
const testAccSecurityGroupConfigUpdate = `
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg-updated"
}
`
