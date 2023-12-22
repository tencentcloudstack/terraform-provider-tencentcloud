package vpc_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudVpcV3RouteEntryBasic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			log.Printf("TF: \n%s", testAccVpcRouteEntryConfig)
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpcRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteEntryExists("tencentcloud_route_table_entry.foo"),

					resource.TestCheckResourceAttr("tencentcloud_route_table_entry.foo", "next_type", "EIP"),
					resource.TestCheckResourceAttr("tencentcloud_route_table_entry.foo", "description", tcacctest.DefaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_route_table_entry.foo", "destination_cidr_block", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_route_table_entry.foo", "next_hub", "0"),
				),
			},
			{
				Config: testAccVpcRouteEntryUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_route_table_entry.foo", "disabled", "true"),
				),
			},
			{
				Config: testAccVpcRouteEntryUpdate2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_route_table_entry.foo", "disabled", "false"),
				),
			},
			{
				ResourceName:            "tencentcloud_route_table_entry.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enabled"},
			},
		},
	})
}

func testAccCheckVpcRouteEntryExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		items := strings.Split(rs.Primary.ID, ".")
		if len(items) != 2 {
			return fmt.Errorf("exist test, entry id [%s] is destroyed, we can not get route table id", rs.Primary.ID)
		}

		routeTableId := items[1]
		entryId, err := strconv.ParseUint(items[0], 10, 64)
		if err != nil {
			return err
		}

		info, has, err := service.DescribeRouteTable(ctx, routeTableId)
		if err != nil {
			return err
		}
		if has == 0 {
			return fmt.Errorf("route table not exists.")
		}
		if has != 1 {
			err = fmt.Errorf("one routeTable id get %d routeTable infos", has)
			return err
		}
		for _, v := range info.EntryInfos() {
			if v.RouteEntryId() == int64(entryId) {
				return nil
			}
		}

		return fmt.Errorf("route table entry not exists.")
	}
}

func testAccCheckVpcRouteEntryDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_route_table_entry" {
			continue
		}

		time.Sleep(5 * time.Second)

		items := strings.Split(rs.Primary.ID, ".")
		if len(items) != 2 {
			return fmt.Errorf("destroy test,entry id be destroyed[%s], we can not get route table id", rs.Primary.ID)
		}

		routeTableId := items[1]
		entryId, err := strconv.ParseUint(items[0], 10, 64)
		if err != nil {
			return err
		}

		info, has, err := service.DescribeRouteTable(ctx, routeTableId)
		if err != nil {
			return err
		}
		if has == 0 {
			return nil
		}
		if has != 1 {
			err = fmt.Errorf("one routeTable id get %d routeTable infos", has)
			return err
		}
		for _, v := range info.EntryInfos() {
			if v.RouteEntryId() == int64(entryId) {
				return fmt.Errorf("route table entry still exists")
			}
		}
	}

	return nil
}

const testAccVpcRouteEntryConfig = tcacctest.DefaultVpcVariable + `
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

resource "tencentcloud_route_table_entry" "foo" {
  route_table_id         = tencentcloud_route_table.foo.id
  destination_cidr_block = "10.0.0.0/24"
  next_type              = "EIP"
  next_hub               = "0"
  description            = var.instance_name
}
`
const testAccVpcRouteEntryUpdate = tcacctest.DefaultVpcVariable + `
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

resource "tencentcloud_route_table_entry" "foo" {
  route_table_id         = tencentcloud_route_table.foo.id
  destination_cidr_block = "10.0.0.0/24"
  next_type              = "EIP"
  next_hub               = "0"
  description            = var.instance_name
  disabled                = true
}
`
const testAccVpcRouteEntryUpdate2 = tcacctest.DefaultVpcVariable + `
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

resource "tencentcloud_route_table_entry" "foo" {
  route_table_id         = tencentcloud_route_table.foo.id
  destination_cidr_block = "10.0.0.0/24"
  next_type              = "EIP"
  next_hub               = "0"
  description            = var.instance_name
  disabled               = false
}
`
