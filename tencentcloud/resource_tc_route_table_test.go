package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudVpcV3RouteTable_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteTableExists("tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", "ci-temp-test-rt"),
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

func TestAccTencentCloudVpcV3RouteTable_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteTableExists("tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", "ci-temp-test-rt"),

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
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", "ci-temp-test-rt-updated"),

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

func testAccCheckVpcRouteTableExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(nil)
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

		return fmt.Errorf("route table not exists.")
	}
}

func testAccCheckVpcRouteTableDestroy(s *terraform.State) error {
	logId := getLogId(nil)
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
		return fmt.Errorf("route table  not delete ok")
	}
	return nil
}

const testAccVpcRouteTableConfig = `
resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test"
    cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "foo" {
   vpc_id = "${tencentcloud_vpc.foo.id}"
   name = "ci-temp-test-rt"
}
`
const testAccVpcRouteTableConfigUpdate = `
resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test"
    cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "foo" {
   vpc_id = "${tencentcloud_vpc.foo.id}"
   name = "ci-temp-test-rt-updated"
}
`
