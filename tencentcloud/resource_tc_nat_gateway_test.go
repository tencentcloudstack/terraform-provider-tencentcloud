package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	vpc "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
)

func TestAccTencentCloudNatGateway_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_nat_gateway.my_nat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "max_concurrent", "3000000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "bandwidth", "500"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "assigned_eip_set.#", "2"),
				),
			},
			{
				Config: testAccNatGatewayConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_nat_gateway.my_nat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "name", "new_name"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "max_concurrent", "10000000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "bandwidth", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "assigned_eip_set.#", "2"),
				),
			},
		},
	})
}

func testAccCheckNatGatewayDestroy(s *terraform.State) error {

	conn := testAccProvider.Meta().(*TencentCloudClient).vpcConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_nat_gateway" {
			continue
		}

		descReq := vpc.NewDescribeNatGatewayRequest()
		descReq.NatId = common.StringPtr(rs.Primary.ID)
		descResp, err := conn.DescribeNatGateway(descReq)

		b, _ := json.Marshal(descResp)

		log.Printf("[DEBUG] conn.DescribeNatGateway response: %s", b)

		if _, ok := err.(*common.APIError); ok {
			return fmt.Errorf("conn.DescribeNatGateway error: %v", err)
		} else if *descResp.TotalCount != 0 {
			return fmt.Errorf("NAT Gateway still exists.")
		}
	}
	return nil
}

const testAccNatGatewayConfig = `
resource "tencentcloud_vpc" "main" {
  name       = "terraform test"
  cidr_block = "10.6.0.0/16"
}

resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}

resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id           = "${tencentcloud_vpc.main.id}"
  name             = "terraform_test"
  max_concurrent   = 3000000
  bandwidth        = 500
  assigned_eip_set = [
    "${tencentcloud_eip.eip_dev_dnat.public_ip}",
    "${tencentcloud_eip.eip_test_dnat.public_ip}",
  ]
}
`
const testAccNatGatewayConfigUpdate = `
resource "tencentcloud_vpc" "main" {
  name       = "terraform test"
  cidr_block = "10.6.0.0/16"
}

resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}

resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}

resource "tencentcloud_eip" "new_eip" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id           = "${tencentcloud_vpc.main.id}"
  name             = "new_name"
  max_concurrent   = 10000000
  bandwidth        = 1000
  assigned_eip_set = [
    "${tencentcloud_eip.eip_dev_dnat.public_ip}",
    "${tencentcloud_eip.new_eip.public_ip}",
  ]
}
`
