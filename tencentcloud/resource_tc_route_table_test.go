package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudVpcV3RouteTableBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteTableExists("tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", defaultInsName),

					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "subnet_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "route_entry_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "create_time"),
				),
			},
			{
				ResourceName:      "tencentcloud_route_table.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudVpcV3RouteTableUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteTableExists("tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", defaultInsName),

					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "subnet_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "route_entry_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "create_time"),
				),
			},
			{
				Config: testAccVpcRouteTableConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteTableExists("tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", defaultInsNameUpdate),

					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "subnet_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "route_entry_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudVpcV3RouteTableWithTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTableConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteTableExists("tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", defaultInsName),

					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "subnet_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "route_entry_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_route_table.foo", "create_time"),

					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "tags.test", "test"),
					resource.TestCheckNoResourceAttr("tencentcloud_route_table.foo", "tags.abc"),
				),
			},
			{
				Config: testAccVpcRouteTableConfigWithTagsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteTableExists("tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", defaultInsName),

					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "tags.abc", "abc"),
					resource.TestCheckNoResourceAttr("tencentcloud_route_table.foo", "tags.test"),
				),
			},
		},
	})
}

func testAccCheckVpcRouteTableExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := service.DescribeRouteTable(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has > 0 {
			return nil
		}

		return fmt.Errorf("routetable %s not exists", rs.Primary.ID)
	}
}

func testAccCheckVpcRouteTableDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_route_table" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, has, err := service.DescribeRouteTable(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has == 0 {
			return nil
		}
		return fmt.Errorf("routetable %s still exists", rs.Primary.ID)
	}

	return nil
}

const testAccVpcRouteTableConfig = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = "${var.instance_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "tencentcloud_route_table" "foo" {
  name   = "${var.instance_name}"
  vpc_id = "${tencentcloud_vpc.foo.id}"
}
`

const testAccVpcRouteTableConfigUpdate = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = "${var.instance_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "tencentcloud_route_table" "foo" {
  name   = "${var.instance_name_update}"
  vpc_id = "${tencentcloud_vpc.foo.id}"
}
`

const testAccVpcRouteTableConfigWithTags = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = "${var.instance_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "tencentcloud_route_table" "foo" {
  name   = "${var.instance_name}"
  vpc_id = "${tencentcloud_vpc.foo.id}"

  tags = {
    "test" = "test"
  }
}
`

const testAccVpcRouteTableConfigWithTagsUpdate = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = "${var.instance_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "tencentcloud_route_table" "foo" {
  name   = "${var.instance_name}"
  vpc_id = "${tencentcloud_vpc.foo.id}"

  tags = {
    "abc" = "abc"
  }
}
`
