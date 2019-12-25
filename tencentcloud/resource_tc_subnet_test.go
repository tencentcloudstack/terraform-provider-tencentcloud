package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudVpcV3SubnetBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetExists("tencentcloud_subnet.subnet"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", defaultSubnetCidr),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", defaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "is_multicast", "false"),
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

func TestAccTencentCloudVpcV3SubnetUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetExists("tencentcloud_subnet.subnet"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", defaultSubnetCidr),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", defaultInsName),
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
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", defaultSubnetCidrLess),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", defaultInsNameUpdate),
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

func TestAccTencentCloudVpcV3SubnetWithTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetExists("tencentcloud_subnet.subnet"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", defaultSubnetCidr),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", defaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "is_multicast", "false"),

					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "availability_zone"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "available_ip_count"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "route_table_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_subnet.subnet", "create_time"),

					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "tags.test", "test"),
					resource.TestCheckNoResourceAttr("tencentcloud_subnet.subnet", "tags.abc"),
				),
			},
			{
				Config: testAccVpcSubnetConfigWithTagsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetExists("tencentcloud_subnet.subnet"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", defaultSubnetCidr),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", defaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "is_multicast", "false"),

					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "tags.abc", "abc"),
					resource.TestCheckNoResourceAttr("tencentcloud_subnet.subnet", "tags.test"),
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

		return fmt.Errorf("subnet %s not exists", rs.Primary.ID)
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
		return fmt.Errorf("subnet %s still exists", rs.Primary.ID)
	}

	return nil
}

const testAccVpcSubnetConfig = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_subnet" "subnet" {
  name              = var.instance_name
  vpc_id            = tencentcloud_vpc.foo.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr
  is_multicast      = false
}
`

const testAccVpcSubnetConfigUpdate = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_route_table" "route_table" {
  name   = var.instance_name
  vpc_id = tencentcloud_vpc.foo.id
}

resource "tencentcloud_subnet" "subnet" {
  name              = var.instance_name_update
  vpc_id            = tencentcloud_vpc.foo.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr_less
  is_multicast      = true
  route_table_id    = tencentcloud_route_table.route_table.id
}
`

const testAccVpcSubnetConfigWithTags = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_subnet" "subnet" {
  name              = var.instance_name
  vpc_id            = tencentcloud_vpc.foo.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr
  is_multicast      = false

  tags = {
    "test" = "test"
  }
}
`

const testAccVpcSubnetConfigWithTagsUpdate = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_subnet" "subnet" {
  name              = var.instance_name
  vpc_id            = tencentcloud_vpc.foo.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr
  is_multicast      = false

  tags = {
    "abc" = "abc"
  }
}
`
