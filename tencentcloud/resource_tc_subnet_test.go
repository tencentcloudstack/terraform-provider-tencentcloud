package tencentcloud

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudSubnet_basic(t *testing.T) {
	var vpcId string
	var subnetId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreSetRegion("ap-guangzhou")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy(&vpcId, &subnetId),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists("tencentcloud_vpc.foo", "tencentcloud_subnet.foo", &vpcId, &subnetId),
					resource.TestCheckResourceAttr("tencentcloud_subnet.foo", "cidr_block", "10.0.11.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.foo", "name", "ci-temp-test-subnet"),
				),
			},
		},
	})
}

func TestAccTencentCloudSubnet_update(t *testing.T) {
	var vpcId string
	var subnetId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubnetDestroy(&vpcId, &subnetId),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists("tencentcloud_vpc.foo", "tencentcloud_subnet.foo", &vpcId, &subnetId),
					resource.TestCheckResourceAttr("tencentcloud_subnet.foo", "cidr_block", "10.0.11.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.foo", "name", "ci-temp-test-subnet"),
				),
			},
			{
				Config: testAccSubnetConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubnetExists("tencentcloud_vpc.foo", "tencentcloud_subnet.foo", &vpcId, &subnetId),
					resource.TestCheckResourceAttr("tencentcloud_subnet.foo", "name", "ci-temp-test-subnet-updated"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.foo", "route_table_id"),
				),
			},
		},
	})
}

func testAccCheckSubnetDestroy(vpcId, subnetId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *vpcId == "" || *subnetId == "" {
			return nil
		}
		conn := testAccProvider.Meta().(*TencentCloudClient).commonConn
		params := map[string]string{
			"Action":   "DescribeSubnet",
			"vpcId":    *vpcId,
			"subnetId": *subnetId,
		}
		response, err := conn.SendRequest("vpc", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Code       int    `json:"code"`
			CodeDesc   string `json:"codeDesc"`
			TotalCount int    `json:"totalCount"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return err
		}
		if jsonresp.CodeDesc != "InvalidSubnet.NotFound" {
			return fmt.Errorf("Subnet still exists.")
		}
		return nil
	}
}

func testAccCheckSubnetExists(vpc, subnet string, vpcId, subnetId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vpcrs, ok := s.RootModule().Resources[vpc]
		if !ok {
			return fmt.Errorf("Not found: %s", vpc)
		}

		if vpcrs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		rs, ok := s.RootModule().Resources[subnet]
		if !ok {
			return fmt.Errorf("Not found: %s", subnet)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No subnet ID is set")
		}

		client := testAccProvider.Meta().(*TencentCloudClient).commonConn
		params := map[string]string{
			"Action":   "DescribeSubnet",
			"vpcId":    vpcrs.Primary.ID,
			"subnetId": rs.Primary.ID,
		}
		response, err := client.SendRequest("vpc", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Code       int    `json:"code"`
			CodeDesc   string `json:"codeDesc"`
			TotalCount int    `json:"totalCount"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return err
		}
		if jsonresp.CodeDesc == "InvalidSubnet.NotFound" {
			return fmt.Errorf("Subnet not found: %s", rs.Primary.ID)
		}
		if jsonresp.Code != 0 {
			return fmt.Errorf("Describe subnet failed")
		}

		*vpcId = vpcrs.Primary.ID
		*subnetId = rs.Primary.ID
		return nil

	}
}

const testAccSubnetConfig = `
resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test"
    cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
   vpc_id = "${tencentcloud_vpc.foo.id}"
   name = "ci-temp-test-subnet"
   cidr_block = "10.0.11.0/24"
   availability_zone = "ap-guangzhou-3"
}
`

const testAccSubnetConfigUpdate = `
resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_route_table" "foo" {
   vpc_id = "${tencentcloud_vpc.foo.id}"
   name = "ci-temp-test-rt"
}
resource "tencentcloud_subnet" "foo" {
  vpc_id = "${tencentcloud_vpc.foo.id}"
  name = "ci-temp-test-subnet-updated"
  cidr_block = "10.0.11.0/24"
  availability_zone = "ap-guangzhou-3"
  route_table_id = "${tencentcloud_route_table.foo.id}"
}
`
