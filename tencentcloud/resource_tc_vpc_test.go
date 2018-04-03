package tencentcloud

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudVpc_basic(t *testing.T) {
	var vpcId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy(&vpcId),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo", &vpcId),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", "10.0.0.0/16"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", "ci-temp-test"),
				),
			},
		},
	})
}

func TestAccTencentCloudVpc_update(t *testing.T) {
	var vpcId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy(&vpcId),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo", &vpcId),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", "10.0.0.0/16"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", "ci-temp-test"),
				),
			},
			{
				Config: testAccVpcConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo", &vpcId),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", "10.0.0.0/16"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", "ci-temp-test-updated"),
				),
			},
		},
	})
}

const testAccVpcConfig = `
resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test"
    cidr_block = "10.0.0.0/16"
}
`

const testAccVpcConfigUpdate = `
resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test-updated"
    cidr_block = "10.0.0.0/16"
}
`

func testAccCheckVpcDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *id == "" {
			return nil
		}
		conn := testAccProvider.Meta().(*TencentCloudClient).commonConn
		params := map[string]string{
			"Action": "DescribeVpcEx",
			"vpcId":  *id,
		}
		response, err := conn.SendRequest("vpc", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Code       int    `json:"code"`
			Message    string `json:"message"`
			TotalCount int    `json:"totalCount"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return err
		}
		if jsonresp.TotalCount != 0 {
			return fmt.Errorf("VPC still exists.")
		}
		return nil
	}
}

func testAccCheckVpcExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		client := testAccProvider.Meta().(*TencentCloudClient).commonConn
		params := map[string]string{
			"Action": "DescribeVpcEx",
			"vpcId":  rs.Primary.ID,
		}
		response, err := client.SendRequest("vpc", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Code       int    `json:"code"`
			Message    string `json:"message"`
			TotalCount int    `json:"totalCount"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return err
		}
		if jsonresp.TotalCount == 0 {
			return fmt.Errorf("VPC not found")
		}

		*id = rs.Primary.ID
		return nil
	}

}
