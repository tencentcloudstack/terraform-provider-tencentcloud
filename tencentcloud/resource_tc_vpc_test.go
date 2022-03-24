package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_vpc", &resource.Sweeper{
		Name: "tencentcloud_vpc",
		F:    testSweepVpcInstance,
	})
}

func testSweepVpcInstance(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sharedClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(*TencentCloudClient)

	vpcService := VpcService{
		client: client.apiV3Conn,
	}

	instances, err := vpcService.DescribeVpcs(ctx, "", "", nil, nil, "", "")
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {
		instanceId := v.vpcId
		instanceName := v.name

		now := time.Now()

		createTime := stringTotime(v.createTime)
		interval := now.Sub(createTime).Minutes()
		if strings.HasPrefix(instanceName, keepResource) || strings.HasPrefix(instanceName, defaultResource) {
			continue
		}
		// less than 30 minute, not delete
		if int64(interval) < 30 {
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", defaultVpcCidr),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", defaultInsName),
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", defaultVpcCidr),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", defaultInsName),
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
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", defaultVpcCidrLess),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", defaultInsNameUpdate),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "is_multicast", "false"),

					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "dns_servers.#"),

					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", fmt.Sprintf("%s.%d", "dns_servers", hashcode.String("119.29.29.29")), "119.29.29.29"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", fmt.Sprintf("%s.%d", "dns_servers", hashcode.String("182.254.116.116")), "182.254.116.116"),
				),
			},
		},
	})
}

func TestAccTencentCloudVpcV3WithTags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", defaultVpcCidr),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", defaultInsName),
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
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", defaultVpcCidr),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", defaultInsName),
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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

const testAccVpcConfig = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}
`

const testAccVpcConfigUpdate = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name_update
  cidr_block = var.vpc_cidr_less
  dns_servers  = ["119.29.29.29", "182.254.116.116"]
  is_multicast = false
}
`

const testAccVpcConfigWithTags = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr

  tags = {
    "test" = "test"
  }
}
`

const testAccVpcConfigWithTagsUpdate = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr

  tags = {
    "abc" = "abc"
  }
}
`
