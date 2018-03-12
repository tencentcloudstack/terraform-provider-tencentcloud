package tencentcloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	vpc "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
)

func TestAccTencentCloudDnat_basic(t *testing.T) {
	var dnatId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnatConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnatExists("tencentcloud_dnat.dev_dnat", &dnatId),
				),
			},
			{
				Config: testAccDnatConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnatExists("tencentcloud_dnat.dev_dnat", &dnatId),
				),
			},
		},
	})
}

func testAccCheckDnatExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s\n", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		_entry, _ := parseDnatId(rs.Primary.ID)

		client := testAccProvider.Meta().(*TencentCloudClient)
		entry, err := client.DescribeDnat(_entry)

		if err != nil {
			return err
		} else if entry == nil {
			return dnatNotFound
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckDnatDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*TencentCloudClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dnat" {
			continue
		}

		_entry := &vpc.DnaptRule{}

		_entry.UniqVpcId = common.StringPtr(rs.Primary.Attributes["vpc_id"])
		_entry.UniqNatId = common.StringPtr(rs.Primary.Attributes["nat_id"])
		_entry.Proto = common.StringPtr(rs.Primary.Attributes["ip_protocol"])
		_entry.Eip = common.StringPtr(rs.Primary.Attributes["external_ip"])
		_external_port, _ := strconv.Atoi(rs.Primary.Attributes["external_port"])
		_entry.Eport = common.IntPtr(_external_port)

		_, err := client.DescribeDnat(_entry)

		if err == dnatNotFound {
			return nil
		} else if err != nil {
			return err
		} else {
			return fmt.Errorf("DNAT still exists.")
		}
	}
	return nil
}

const testAccDnatConfig = `
data "tencentcloud_availability_zones" "my_favorate_zones" {}

data "tencentcloud_image" "my_favorate_image" {
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

# Create VPC and Subnet
resource "tencentcloud_vpc" "main" {
  name       = "terraform test"
  cidr_block = "10.6.0.0/16"
}
resource "tencentcloud_subnet" "main_subnet" {
  vpc_id            = "${tencentcloud_vpc.main.id}"
  name              = "terraform test subnet"
  cidr_block        = "10.6.7.0/24"
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
}

# Create EIP 
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}

# Create NAT Gateway
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id           = "${tencentcloud_vpc.main.id}"
  name             = "terraform test"
  max_concurrent   = 3000000
  bandwidth        = 500
  assigned_eip_set = [
    "${tencentcloud_eip.eip_dev_dnat.public_ip}",
    "${tencentcloud_eip.eip_test_dnat.public_ip}",
  ]
}

# Create CVM
resource "tencentcloud_instance" "foo" {
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
  vpc_id            = "${tencentcloud_vpc.main.id}"
  subnet_id         = "${tencentcloud_subnet.main_subnet.id}"
}

# Add DNAT Entry
resource "tencentcloud_dnat" "dev_dnat" {
  vpc_id       = "${tencentcloud_nat_gateway.my_nat.vpc_id}"
  nat_id       = "${tencentcloud_nat_gateway.my_nat.id}"
  protocol     = "tcp"
  elastic_ip   = "${tencentcloud_eip.eip_dev_dnat.public_ip}"
  elastic_port = "80"
  private_ip   = "${tencentcloud_instance.foo.private_ip}"
  private_port = "9001"
}
`
const testAccDnatConfigUpdate = `
data "tencentcloud_availability_zones" "my_favorate_zones" {}

data "tencentcloud_image" "my_favorate_image" {
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

# Create VPC and Subnet
resource "tencentcloud_vpc" "main" {
  name       = "terraform test"
  cidr_block = "10.6.0.0/16"
}
resource "tencentcloud_subnet" "main_subnet" {
  vpc_id            = "${tencentcloud_vpc.main.id}"
  name              = "terraform test subnet"
  cidr_block        = "10.6.7.0/24"
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
}

# Create EIP 
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}

# Create NAT Gateway
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id           = "${tencentcloud_vpc.main.id}"
  name             = "terraform test"
  max_concurrent   = 3000000
  bandwidth        = 500
  assigned_eip_set = [
    "${tencentcloud_eip.eip_dev_dnat.public_ip}",
    "${tencentcloud_eip.eip_test_dnat.public_ip}",
  ]
}

# Create CVM
resource "tencentcloud_instance" "foo" {
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
  vpc_id            = "${tencentcloud_vpc.main.id}"
  subnet_id         = "${tencentcloud_subnet.main_subnet.id}"
}

# Update DNAT Entry
resource "tencentcloud_dnat" "dev_dnat" {
  vpc_id       = "${tencentcloud_nat_gateway.my_nat.vpc_id}"
  nat_id       = "${tencentcloud_nat_gateway.my_nat.id}"
  protocol     = "udp"
  elastic_ip   = "${tencentcloud_eip.eip_test_dnat.public_ip}"
  elastic_port = "8080"
  private_ip   = "${tencentcloud_instance.foo.private_ip}"
  private_port = "9002"
}
`
