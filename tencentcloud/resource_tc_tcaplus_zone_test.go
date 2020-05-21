package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testTcaplusZoneResourceName = "tencentcloud_tcaplus_zone"
var testTcaplusZoneResourceNameResourceKey = testTcaplusZoneResourceName + ".test_zone"

func TestAccTencentCloudTcaplusZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcaplusZone,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusZoneExists(testTcaplusZoneResourceNameResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusZoneResourceNameResourceKey, "total_size"),
					resource.TestCheckResourceAttrSet(testTcaplusZoneResourceNameResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testTcaplusZoneResourceNameResourceKey, "zone_name", "tf_test_zone_name_guagua"),
					resource.TestCheckResourceAttr(testTcaplusZoneResourceNameResourceKey, "table_count", "0"),
				),
			},
			{
				Config: testAccTcaplusZoneUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTcaplusZoneExists(testTcaplusZoneResourceNameResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusZoneResourceNameResourceKey, "total_size"),
					resource.TestCheckResourceAttrSet(testTcaplusZoneResourceNameResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testTcaplusZoneResourceNameResourceKey, "zone_name", "tf_test_zone_name_guagua_2"),
					resource.TestCheckResourceAttr(testTcaplusZoneResourceNameResourceKey, "table_count", "0"),
				),
			},
		},
	})
}
func testAccCheckTcaplusZoneDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTcaplusZoneResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeZone(ctx, rs.Primary.Attributes["app_id"], rs.Primary.ID)

		if err != nil {
			_, has, err = service.DescribeZone(ctx, rs.Primary.Attributes["app_id"], rs.Primary.ID)
		}

		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("delete tcaplus zone %s fail, still on server", rs.Primary.ID)
	}
	return nil
}

func testAccCheckTcaplusZoneExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeZone(ctx, rs.Primary.Attributes["app_id"], rs.Primary.ID)

		if err != nil {
			_, has, err = service.DescribeZone(ctx, rs.Primary.Attributes["app_id"], rs.Primary.ID)
		}

		if err != nil {
			return err
		}
		if has {
			return nil
		}
		return fmt.Errorf("tcaplus zone %s not found on server", rs.Primary.ID)
	}
}

const testAccTcaplusZoneBasic = `variable "availability_zone" {
  default = "ap-shanghai-2"
}
variable "instance_name" {
  default = "` + defaultInsName + `"
}
variable "vpc_cidr" {
  default = "` + defaultVpcCidr + `"
}
variable "subnet_cidr" {
  default = "` + defaultSubnetCidr + `"
}

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
resource "tencentcloud_tcaplus_application" "test_app" {
  idl_type                 = "PROTO"
  app_name                 = "tf_guagua_zone"
  vpc_id                   = tencentcloud_vpc.foo.id
  subnet_id                = tencentcloud_subnet.subnet.id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}`

const testAccTcaplusZone = testAccTcaplusZoneBasic + `
resource "tencentcloud_tcaplus_zone" "test_zone" {
  app_id    = tencentcloud_tcaplus_application.test_app.id
  zone_name = "tf_test_zone_name_guagua"
}
`
const testAccTcaplusZoneUpdate = testAccTcaplusZoneBasic + `
resource "tencentcloud_tcaplus_zone" "test_zone" {
  app_id    = tencentcloud_tcaplus_application.test_app.id
  zone_name = "tf_test_zone_name_guagua_2"
}
`
