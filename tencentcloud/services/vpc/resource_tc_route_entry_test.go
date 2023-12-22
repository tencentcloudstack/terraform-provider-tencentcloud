package vpc_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudVpcV2RouteEntryBasic(t *testing.T) {
	t.Parallel()
	var reId string
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		route, ok := svcvpc.RouteIdDecode(*id)
		if !ok {
			return fmt.Errorf("tencentcloud_route_entry read error, id decode faild, id:%v", id)
		}

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			info, has, e := service.DescribeRouteTable(ctx, route["routeTableId"])
			if e != nil {
				return tccommon.RetryError(e)
			}
			if has == 0 {
				return nil
			}
			if has != 1 {
				e = fmt.Errorf("one routeTable id get %d routeTable infos", has)
				return resource.NonRetryableError(e)
			}
			for _, v := range info.EntryInfos() {
				var nextType string
				var nextTypeId string
				for kk, vv := range svcvpc.RouteTypeNewMap {
					if vv == v.NextType() {
						nextType = kk
					}
				}
				if _, ok := svcvpc.RouteTypeApiMap[nextType]; ok {
					nextTypeId = fmt.Sprintf("%d", svcvpc.RouteTypeApiMap[nextType])
				}
				if v.DestinationCidr() == route["destinationCidrBlock"] &&
					nextTypeId == route["nextType"] &&
					v.NextBub() == route["nextHub"] &&
					v.Description() == route["description"] {
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

		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		route, ok := svcvpc.RouteIdDecode(rs.Primary.ID)
		if !ok {
			return fmt.Errorf("tencentcloud_route_entry read error, id decode faild, id:%v", rs.Primary.ID)
		}

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			info, has, e := service.DescribeRouteTable(ctx, route["routeTableId"])
			if e != nil {
				return tccommon.RetryError(e)
			}
			if has == 0 {
				return resource.NonRetryableError(fmt.Errorf("Route entry not found: %s", rs.Primary.ID))
			}
			if has != 1 {
				e = fmt.Errorf("one routeTable id get %d routeTable infos", has)
				return resource.NonRetryableError(e)
			}
			for _, v := range info.EntryInfos() {
				var nextType string
				var nextTypeId string
				for kk, vv := range svcvpc.RouteTypeNewMap {
					if vv == v.NextType() {
						nextType = kk
					}
				}
				if _, ok := svcvpc.RouteTypeApiMap[nextType]; ok {
					nextTypeId = fmt.Sprintf("%d", svcvpc.RouteTypeApiMap[nextType])
				}
				if v.DestinationCidr() == route["destinationCidrBlock"] &&
					nextTypeId == route["nextType"] &&
					v.NextBub() == route["nextHub"] &&
					v.Description() == route["description"] {
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

const testAccVpcRouteEntryV2Config = tcacctest.DefaultVpcVariable + `
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
