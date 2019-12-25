package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudVpcV2RouteEntryBasic(t *testing.T) {
	var reId string
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteEntryDestroy(&reId),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteEntryV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteEntryExists("tencentcloud_route_entry.foo", &reId),
					resource.TestCheckResourceAttr("tencentcloud_route_entry.foo", "cidr_block", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_route_entry.foo", "next_type", "eip"),
				),
			},
		},
	})
}

func testAccCheckRouteEntryDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)
		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		route, ok := routeIdDecode(*id)
		if !ok {
			return fmt.Errorf("tencentcloud_route_entry read error, id decode faild, id:%v", id)
		}

		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, has, e := service.DescribeRouteTable(ctx, route["routeTableId"])
			if e != nil {
				return retryError(e)
			}
			if has == 0 {
				return nil
			}
			if has != 1 {
				e = fmt.Errorf("one routeTable id get %d routeTable infos", has)
				return resource.NonRetryableError(e)
			}
			for _, v := range info.entryInfos {
				var nextType string
				var nextTypeId string
				for kk, vv := range routeTypeNewMap {
					if vv == v.nextType {
						nextType = kk
					}
				}
				if _, ok := routeTypeApiMap[nextType]; ok {
					nextTypeId = fmt.Sprintf("%d", routeTypeApiMap[nextType])
				}
				if v.destinationCidr == route["destinationCidrBlock"] &&
					nextTypeId == route["nextType"] &&
					v.nextBub == route["nextHub"] &&
					v.description == route["description"] {
					return resource.NonRetryableError(fmt.Errorf("Route entry still exists: %s", *id))
				}
			}

			return nil
		})
		return err
	}
}

func testAccCheckRouteEntryExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No route entry ID is set")
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)
		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		route, ok := routeIdDecode(rs.Primary.ID)
		if !ok {
			return fmt.Errorf("tencentcloud_route_entry read error, id decode faild, id:%v", rs.Primary.ID)
		}

		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, has, e := service.DescribeRouteTable(ctx, route["routeTableId"])
			if e != nil {
				return retryError(e)
			}
			if has == 0 {
				return resource.NonRetryableError(fmt.Errorf("Route entry not found: %s", rs.Primary.ID))
			}
			if has != 1 {
				e = fmt.Errorf("one routeTable id get %d routeTable infos", has)
				return resource.NonRetryableError(e)
			}
			for _, v := range info.entryInfos {
				var nextType string
				var nextTypeId string
				for kk, vv := range routeTypeNewMap {
					if vv == v.nextType {
						nextType = kk
					}
				}
				if _, ok := routeTypeApiMap[nextType]; ok {
					nextTypeId = fmt.Sprintf("%d", routeTypeApiMap[nextType])
				}
				if v.destinationCidr == route["destinationCidrBlock"] &&
					nextTypeId == route["nextType"] &&
					v.nextBub == route["nextHub"] &&
					v.description == route["description"] {
					return nil
				}
			}

			return resource.NonRetryableError(fmt.Errorf("Route entry not found: %s", rs.Primary.ID))
		})
		if err != nil {
			return err
		}

		*id = rs.Primary.ID
		return nil
	}
}

const testAccVpcRouteEntryV2Config = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_route_table" "foo" {
  name   = var.instance_name
  vpc_id = tencentcloud_vpc.foo.id
}

resource "tencentcloud_subnet" "foo" {
  name              = var.instance_name
  vpc_id            = tencentcloud_vpc.foo.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr
  is_multicast      = false
  route_table_id    = tencentcloud_route_table.foo.id
}

resource "tencentcloud_route_entry" "foo" {
  vpc_id        	= tencentcloud_vpc.foo.id
  route_table_id 	= tencentcloud_route_table.foo.id
  cidr_block 		= "10.0.0.0/24"
  next_type 		= "eip"
  next_hub  		= "0"
}
`
