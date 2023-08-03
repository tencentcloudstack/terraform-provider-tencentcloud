package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcIpv6EniAddressResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		//CheckDestroy: testAccCheckVpcIpv6EniAddressDestroy,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcIpv6EniAddress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcIpv6EniAddressExists("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "network_interface_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "ipv6_addresses.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "ipv6_addresses.0.address"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_ipv6_eni_address.ipv6_eni_address", "ipv6_addresses.0.description"),
				),
			},
		},
	})
}

func testAccCheckVpcIpv6EniAddressDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpc_ipv6_eni_address" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		vpcId := idSplit[0]
		//networkInterfaceId := idSplit[1]
		address := idSplit[2]

		ipv6EniAddress, err := service.DescribeVpcIpv6EniAddressById(ctx, vpcId, address)
		if err != nil {
			return err
		}
		if ipv6EniAddress != nil {
			return fmt.Errorf("VpcIpv6EniAddress Resource %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckVpcIpv6EniAddressExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		vpcId := idSplit[0]
		//networkInterfaceId := idSplit[1]
		address := idSplit[2]

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		ipv6EniAddress, err := service.DescribeVpcIpv6EniAddressById(ctx, vpcId, address)
		if err != nil {
			return err
		}
		if ipv6EniAddress == nil {
			return fmt.Errorf("VpcIpv6EniAddress Resource %s is not found", rs.Primary.ID)
		}
		return nil
	}
}

const testAccVpcIpv6EniAddress = `

resource "tencentcloud_vpc" "foo" {
  name       = "ci-test-eni-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = tencentcloud_vpc.foo.id
  subnet_id   = tencentcloud_subnet.foo.id
  description = "eni desc"
  ipv4_count  = 1
}

resource "tencentcloud_subnet" "foo" {
  availability_zone = "ap-guangzhou-3"
  name              = "ci-test-eni-subnet"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_vpc_ipv6_subnet_cidr_block" "ipv6_subnet_cidr_block" {
  vpc_id = tencentcloud_vpc.foo.id
  ipv6_subnet_cidr_blocks  {
    subnet_id = tencentcloud_subnet.foo.id
    ipv6_cidr_block = replace(tencentcloud_vpc_ipv6_cidr_block.ipv6_cidr_block.ipv6_cidr_block,"/56","/64")
  }
  depends_on = [tencentcloud_eni.foo]
}

resource "tencentcloud_vpc_ipv6_cidr_block" "ipv6_cidr_block" {
  vpc_id = tencentcloud_vpc.foo.id
}

resource "tencentcloud_vpc_ipv6_eni_address" "ipv6_eni_address" {
  vpc_id      = tencentcloud_vpc.foo.id
  network_interface_id = tencentcloud_eni.foo.id
  ipv6_addresses {
    address = replace(tencentcloud_vpc_ipv6_cidr_block.ipv6_cidr_block.ipv6_cidr_block,":/56","0:99cc:65a7:3017")
    description = "test123"
  }
  depends_on = [tencentcloud_vpc_ipv6_subnet_cidr_block.ipv6_subnet_cidr_block]
}
`
