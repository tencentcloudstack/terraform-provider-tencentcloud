package tencentcloud

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudRouteTable_basic(t *testing.T) {
	var vpcId string
	var rtId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableDestroy(&vpcId, &rtId),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists("tencentcloud_vpc.foo", "tencentcloud_route_table.foo", &vpcId, &rtId),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", "ci-temp-test-rt"),
				),
			},
		},
	})
}

func TestAccTencentCloudRouteTable_update(t *testing.T) {
	var vpcId string
	var rtId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableDestroy(&vpcId, &rtId),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists("tencentcloud_vpc.foo", "tencentcloud_route_table.foo", &vpcId, &rtId),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", "ci-temp-test-rt"),
				),
			},
			{
				Config: testAccRouteTableConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists("tencentcloud_vpc.foo", "tencentcloud_route_table.foo", &vpcId, &rtId),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", "ci-temp-test-rt-updated"),
				),
			},
		},
	})
}

func testAccCheckRouteTableDestroy(vpcId, rtId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*TencentCloudClient).commonConn
		params := map[string]string{
			"Action":       "DescribeRouteTable",
			"vpcId":        *vpcId,
			"routeTableId": *rtId,
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
		if jsonresp.TotalCount != 0 {
			return fmt.Errorf("Route table still exists.")
		}
		return nil
	}
}

func testAccCheckRouteTableExists(vpc, rt string, vpcId, rtId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vpcrs, ok := s.RootModule().Resources[vpc]
		if !ok {
			return fmt.Errorf("Not found: %s", vpc)
		}

		if vpcrs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		rs, ok := s.RootModule().Resources[rt]
		if !ok {
			return fmt.Errorf("Not found: %s", rt)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No route table ID is set")
		}

		client := testAccProvider.Meta().(*TencentCloudClient).commonConn
		params := map[string]string{
			"Action":       "DescribeRouteTable",
			"vpcId":        vpcrs.Primary.ID,
			"routeTableId": rs.Primary.ID,
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
		if jsonresp.TotalCount == 0 {
			return fmt.Errorf("Route table not found: %s", rs.Primary.ID)
		}

		*vpcId = vpcrs.Primary.ID
		*rtId = rs.Primary.ID
		return nil
	}
}

const testAccRouteTableConfig = `
resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test"
    cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "foo" {
   vpc_id = "${tencentcloud_vpc.foo.id}"
   name = "ci-temp-test-rt"
}
`
const testAccRouteTableConfigUpdate = `
resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test"
    cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "foo" {
   vpc_id = "${tencentcloud_vpc.foo.id}"
   name = "ci-temp-test-rt-updated"
}
`
