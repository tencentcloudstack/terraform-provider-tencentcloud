package vpc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudVpcV3RouteTableBasic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpcRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteTableExists("tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", tcacctest.DefaultInsName),

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
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpcRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteTableExists("tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", tcacctest.DefaultInsName),

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
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", tcacctest.DefaultInsNameUpdate),

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
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpcRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTableConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteTableExists("tencentcloud_route_table.foo"),
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", tcacctest.DefaultInsName),

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
					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "name", tcacctest.DefaultInsName),

					resource.TestCheckResourceAttr("tencentcloud_route_table.foo", "tags.abc", "abc"),
					resource.TestCheckNoResourceAttr("tencentcloud_route_table.foo", "tags.test"),
				),
			},
		},
	})
}

func testAccCheckVpcRouteTableExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccVpcRouteTableConfig = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_route_table" "foo" {
  name   = var.instance_name
  vpc_id = tencentcloud_vpc.foo.id
}
`

const testAccVpcRouteTableConfigUpdate = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_route_table" "foo" {
  name   = var.instance_name_update
  vpc_id = tencentcloud_vpc.foo.id
}
`

const testAccVpcRouteTableConfigWithTags = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_route_table" "foo" {
  name   = var.instance_name
  vpc_id = tencentcloud_vpc.foo.id

  tags = {
    "test" = "test"
  }
}
`

const testAccVpcRouteTableConfigWithTagsUpdate = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_route_table" "foo" {
  name   = var.instance_name
  vpc_id = tencentcloud_vpc.foo.id

  tags = {
    "abc" = "abc"
  }
}
`
