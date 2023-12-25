package vpc_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_subnet", &resource.Sweeper{
		Name: "tencentcloud_subnet",
		F:    testSweepSubnet,
	})
}

func testSweepSubnet(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta)

	vpcService := svcvpc.NewVpcService(client.GetAPIV3Conn())

	instances, err := vpcService.DescribeSubnets(ctx, "", "", "", "",
		nil, nil, nil, "", "")
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {

		instanceId := v.SubnetId()
		instanceName := v.Name()
		now := time.Now()
		createTime := tccommon.StringToTime(v.CreateTime())
		interval := now.Sub(createTime).Minutes()

		if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
			continue
		}

		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = vpcService.DeleteSubnet(ctx, instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudVpcV3SubnetBasic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpcSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetExists("tencentcloud_subnet.subnet"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", tcacctest.DefaultSubnetCidr),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", tcacctest.DefaultInsName),
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
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpcSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetExists("tencentcloud_subnet.subnet"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", tcacctest.DefaultSubnetCidr),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", tcacctest.DefaultInsName),
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
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", tcacctest.DefaultSubnetCidrLess),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", tcacctest.DefaultInsNameUpdate),
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
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpcSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetExists("tencentcloud_subnet.subnet"),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", tcacctest.DefaultSubnetCidr),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", tcacctest.DefaultInsName),
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
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "cidr_block", tcacctest.DefaultSubnetCidr),
					resource.TestCheckResourceAttr("tencentcloud_subnet.subnet", "name", tcacctest.DefaultInsName),
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, has, err := service.DescribeSubnet(ctx, rs.Primary.ID, nil, "", "")
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_subnet" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, has, err := service.DescribeSubnet(ctx, rs.Primary.ID, nil, "", "")
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

const testAccVpcSubnetConfig = tcacctest.DefaultVpcVariable + `
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

const testAccVpcSubnetConfigUpdate = tcacctest.DefaultVpcVariable + `
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

const testAccVpcSubnetConfigWithTags = tcacctest.DefaultVpcVariable + `
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

const testAccVpcSubnetConfigWithTagsUpdate = tcacctest.DefaultVpcVariable + `
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
