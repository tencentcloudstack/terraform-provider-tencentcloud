package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudReserveIpAddressesResource_SetIpAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
			tcacctest.AccStepSetRegion(t, "ap-qingyuan")
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccReserveIpAddress,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_reserve_ip_address.reserve_ip", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_reserve_ip_address.reserve_ip", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_reserve_ip_address.reserve_ip", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_reserve_ip_address.reserve_ip", "ip_address", "10.0.0.13"),
					resource.TestCheckResourceAttr("tencentcloud_reserve_ip_address.reserve_ip", "name", "reserve-ip-tf"),
					resource.TestCheckResourceAttr("tencentcloud_reserve_ip_address.reserve_ip", "description", "description"),
					resource.TestCheckResourceAttr("tencentcloud_reserve_ip_address.reserve_ip", "tags.test1", "test1"),
				),
			},
			{
				Config: testAccReserveIpAddressUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_reserve_ip_address.reserve_ip", "id"),
					resource.TestCheckResourceAttr("tencentcloud_reserve_ip_address.reserve_ip", "name", "reserve-ip-tf-update"),
					resource.TestCheckResourceAttr("tencentcloud_reserve_ip_address.reserve_ip", "description", "description-update"),
					resource.TestCheckResourceAttr("tencentcloud_reserve_ip_address.reserve_ip", "tags.test1", "test1"),
					resource.TestCheckResourceAttr("tencentcloud_reserve_ip_address.reserve_ip", "tags.test2", "test2"),
				),
			},
			{
				ResourceName:            "tencentcloud_reserve_ip_address.reserve_ip",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"subnet_id"},
			}},
	})
}

func TestAccTencentCloudReserveIpAddressesResource_NotSetIpAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
			tcacctest.AccStepSetRegion(t, "ap-qingyuan")
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccReserveIpAddressNotSetIpAddress,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_reserve_ip_address.reserve_ip", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_reserve_ip_address.reserve_ip", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_reserve_ip_address.reserve_ip", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_reserve_ip_address.reserve_ip", "ip_address"),
					resource.TestCheckResourceAttr("tencentcloud_reserve_ip_address.reserve_ip", "name", "reserve-ip-tf"),
					resource.TestCheckResourceAttr("tencentcloud_reserve_ip_address.reserve_ip", "description", "description"),
					resource.TestCheckResourceAttr("tencentcloud_reserve_ip_address.reserve_ip", "tags.test1", "test1"),
				),
			},
			{
				ResourceName:            "tencentcloud_reserve_ip_address.reserve_ip",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"subnet_id"},
			}},
	})
}

const testAccReserveIpAddress = `
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-reserve-ip"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-reserve-ip"
  cidr_block        = "10.0.0.0/16"
  availability_zone = "ap-qingyuan-1"
  is_multicast      = false
}

resource "tencentcloud_reserve_ip_address" "reserve_ip" {
  vpc_id = tencentcloud_vpc.vpc.id
  subnet_id = tencentcloud_subnet.subnet.id
  ip_address = "10.0.0.13"
  name = "reserve-ip-tf"
  description = "description"
  tags ={
    "test1" = "test1"
  }
}
`

const testAccReserveIpAddressUpdate = `
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-reserve-ip"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-reserve-ip"
  cidr_block        = "10.0.0.0/16"
  availability_zone = "ap-qingyuan-1"
  is_multicast      = false
}

resource "tencentcloud_reserve_ip_address" "reserve_ip" {
  vpc_id = tencentcloud_vpc.vpc.id
  subnet_id = tencentcloud_subnet.subnet.id
  ip_address = "10.0.0.13"
  name = "reserve-ip-tf-update"
  description = "description-update"
  tags ={
    "test1" = "test1"
    "test2" = "test2"
  }
}
`

const testAccReserveIpAddressNotSetIpAddress = `
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-reserve-ip"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-reserve-ip"
  cidr_block        = "10.0.0.0/16"
  availability_zone = "ap-qingyuan-1"
  is_multicast      = false
}

resource "tencentcloud_reserve_ip_address" "reserve_ip" {
  vpc_id = tencentcloud_vpc.vpc.id
  subnet_id = tencentcloud_subnet.subnet.id
  name = "reserve-ip-tf"
  description = "description"
  tags ={
    "test1" = "test1"
  }
}
`
