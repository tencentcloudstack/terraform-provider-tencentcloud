package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudRouteEntry_basic(t *testing.T) {
	var reId string
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreSetRegion("ap-guangzhou")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteEntryDestroy(&reId),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteEntryExists("tencentcloud_route_entry.instance", &reId),
					resource.TestCheckResourceAttr("tencentcloud_route_entry.instance", "next_type", "instance"),
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

const testAccRouteEntryConfig = `
data "tencentcloud_image" "foo" {
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

resource "tencentcloud_vpc" "foo" {
  name       = "ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  vpc_id            = "${tencentcloud_vpc.foo.id}"
  name              = "terraform test subnet"
  cidr_block        = "10.0.12.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_route_table" "foo" {
  vpc_id = "${tencentcloud_vpc.foo.id}"
  name   = "ci-temp-test-rt"
}

resource "tencentcloud_instance" "foo" {
  availability_zone = "ap-guangzhou-3"
  image_id          = "${data.tencentcloud_image.foo.image_id}"
  vpc_id            = "${tencentcloud_vpc.foo.id}"
  subnet_id         = "${tencentcloud_subnet.foo.id}"
  instance_type     = "S2.SMALL2"
  system_disk_type  = "CLOUD_SSD"
}

resource "tencentcloud_route_entry" "instance" {
  vpc_id         = "${tencentcloud_vpc.foo.id}"
  route_table_id = "${tencentcloud_route_table.foo.id}"
  cidr_block     = "10.4.4.0/24"
  next_type      = "instance"
  next_hub       = "${tencentcloud_instance.foo.private_ip}"
}
`
