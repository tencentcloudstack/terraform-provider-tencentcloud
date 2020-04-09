package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testTcaplusApplicationResourceName = "tencentcloud_tcaplus_application"
var testTcaplusApplicationResourceKey = testTcaplusApplicationResourceName + ".test_app"

func TestAccTencentCloudTcaplusApplicationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcaplusApplication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusApplicationExists(testTcaplusApplicationResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "network_type"),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "password_status"),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "api_access_id"),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "api_access_ip"),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "api_access_port"),

					resource.TestCheckResourceAttr(testTcaplusApplicationResourceKey, "idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusApplicationResourceKey, "app_name", "tf_te1_guagua"),
					resource.TestCheckResourceAttr(testTcaplusApplicationResourceKey, "password", "1qaA2k1wgvfa3ZZZ"),
					resource.TestCheckResourceAttr(testTcaplusApplicationResourceKey, "old_password_expire_last", "3600"),
				),
			},
			{
				ResourceName:            testTcaplusApplicationResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"old_password_expire_last", "password"},
			},

			{
				Config: testAccTcaplusApplicationUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcaplusApplicationExists(testTcaplusApplicationResourceKey),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "network_type"),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "password_status"),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "api_access_id"),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "api_access_ip"),
					resource.TestCheckResourceAttrSet(testTcaplusApplicationResourceKey, "api_access_port"),

					resource.TestCheckResourceAttr(testTcaplusApplicationResourceKey, "idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testTcaplusApplicationResourceKey, "app_name", "tf_te1_guagua_2"),
					resource.TestCheckResourceAttr(testTcaplusApplicationResourceKey, "password", "aQQ2345677888"),
					resource.TestCheckResourceAttr(testTcaplusApplicationResourceKey, "old_password_expire_last", "300"),
				),
			},
		},
	})
}

func testAccCheckTcaplusApplicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTcaplusApplicationResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeApp(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeApp(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete tcaplus application %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTcaplusApplicationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		service := TcaplusService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeApp(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeApp(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("tcaplus application %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccTcaplusApplicationBaic = `
variable "availability_zone" {
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
}`

const testAccTcaplusApplication string = testAccTcaplusApplicationBaic + `
resource "tencentcloud_tcaplus_application" "test_app" {
  idl_type                 = "PROTO"
  app_name                 = "tf_te1_guagua"
  vpc_id                   = tencentcloud_vpc.foo.id
  subnet_id                = tencentcloud_subnet.subnet.id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
`
const testAccTcaplusApplicationUpdate string = testAccTcaplusApplicationBaic + `
resource "tencentcloud_tcaplus_application" "test_app" {
  idl_type                 = "PROTO"
  app_name                 = "tf_te1_guagua_2"
  vpc_id                   = tencentcloud_vpc.foo.id
  subnet_id                = tencentcloud_subnet.subnet.id
  password                 = "aQQ2345677888"
  old_password_expire_last = 300
}
`
