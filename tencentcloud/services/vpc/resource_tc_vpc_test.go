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
	resource.AddTestSweepers("tencentcloud_vpc", &resource.Sweeper{
		Name: "tencentcloud_vpc",
		F:    testSweepVpcInstance,
	})
}

func testSweepVpcInstance(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta)

	vpcService := svcvpc.NewVpcService(client.GetAPIV3Conn())

	instances, err := vpcService.DescribeVpcs(ctx, "", "", nil, nil, "", "")
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {
		instanceId := v.VpcId()
		instanceName := v.Name()

		now := time.Now()

		createTime := tccommon.StringToTime(v.CreateTime())
		interval := now.Sub(createTime).Minutes()
		if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
			continue
		}
		// less than 30 minute, not delete
		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = vpcService.DeleteVpc(ctx, instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudVpcV3Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", tcacctest.DefaultVpcCidr),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", tcacctest.DefaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "is_multicast", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "default_route_table_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudVpcV3Update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", tcacctest.DefaultVpcCidr),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", tcacctest.DefaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "is_multicast", "true"),

					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "dns_servers.#"),
				),
			},
			{
				Config: testAccVpcConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", tcacctest.DefaultVpcCidrLess),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", tcacctest.DefaultInsNameUpdate),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "is_multicast", "false"),

					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "dns_servers.#"),

					resource.TestCheckTypeSetElemAttr("tencentcloud_vpc.foo", "dns_servers.*", "119.29.29.29"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_vpc.foo", "dns_servers.*", "182.254.116.116"),
				),
			},
		},
	})
}

func TestAccTencentCloudVpcV3WithTags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", tcacctest.DefaultVpcCidr),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", tcacctest.DefaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "is_multicast", "true"),

					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "dns_servers.#"),

					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "tags.test", "test"),
					resource.TestCheckNoResourceAttr("tencentcloud_vpc.foo", "tags.abc"),
				),
			},
			{
				Config: testAccVpcConfigWithTagsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", tcacctest.DefaultVpcCidr),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", tcacctest.DefaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "is_multicast", "true"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "tags.abc", "abc"),
					resource.TestCheckNoResourceAttr("tencentcloud_vpc.foo", "tags.test"),
				),
			},
		},
	})
}

func testAccCheckVpcExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, has, err := service.DescribeVpc(ctx, rs.Primary.ID, "", "")
		if err != nil {
			return err
		}
		if has > 0 {
			return nil
		}

		return fmt.Errorf("vpc %s not exists", rs.Primary.ID)
	}
}

func testAccCheckVpcDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpc" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, has, err := service.DescribeVpc(ctx, rs.Primary.ID, "", "")
		if err != nil {
			return err
		}
		if has == 0 {
			return nil
		}
		return fmt.Errorf("vpc %s still exists", rs.Primary.ID)
	}

	return nil
}

const testAccVpcConfig = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}
`

const testAccVpcConfigUpdate = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name_update
  cidr_block = var.vpc_cidr_less
  dns_servers  = ["119.29.29.29", "182.254.116.116"]
  is_multicast = false
}
`

const testAccVpcConfigWithTags = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr

  tags = {
    "test" = "test"
  }
}
`

const testAccVpcConfigWithTagsUpdate = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr

  tags = {
    "abc" = "abc"
  }
}
`
