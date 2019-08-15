package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudVpcV3Subnet_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetExists("tencentcloud_subnet.subnet"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", "10.0.20.0/28"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", "guagua-ci-temp-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_subnet.subnet",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudVpcV3Subnet_update(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetExists("tencentcloud_subnet.subnet"),

					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", "10.0.20.0/28"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", "guagua-ci-temp-test"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "is_multicast", "false"),

					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "availability_zone"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "available_ip_count"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "route_table_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "create_time"),
				),
			},
			{
				Config: testAccVpcSubnetConfigUpdate,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckVpcSubnetExists("tencentcloud_subnet.subnet"),

					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", "10.0.20.0/28"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", "ci-temp-test-subnet-updated"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "is_multicast", "true"),

					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "availability_zone"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "available_ip_count"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "route_table_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "create_time"),
				),
			},
		},
	})
}
func testAccCheckVpcSubnetExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := service.DescribeSubnet(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has > 0 {
			return nil
		}

		return fmt.Errorf("subnet not exists.")
	}
}

func testAccCheckVpcSubnetDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_subnet" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, has, err := service.DescribeSubnet(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has == 0 {
			return nil
		}
		return fmt.Errorf("subnet not delete ok")
	}
	return nil
}

const testAccVpcSubnetConfig = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "guagua-ci-temp-test"
  cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
  availability_zone = "${var.availability_zone}"
  name              = "guagua-ci-temp-test"
  vpc_id            = "${tencentcloud_vpc.foo.id}"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}
`

const testAccVpcSubnetConfigUpdate = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "guagua-ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "route_table" {
  vpc_id = "${tencentcloud_vpc.foo.id}"
  name   = "ci-temp-test-rt"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = "${var.availability_zone}"
  name              = "ci-temp-test-subnet-updated"
  vpc_id            = "${tencentcloud_vpc.foo.id}"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = true
  route_table_id    = "${tencentcloud_route_table.route_table.id}"
}
`
