package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCamRole_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRole_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamRoleExists("tencentcloud_cam_role.role_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role.role_basic", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role.role_basic", "document"),
				),
			}, {
				Config: testAccCamRole_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamRoleExists("tencentcloud_cam_role.role_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role.role_basic", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role.role_basic", "document"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_role.role_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamRoleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_role" {
			continue
		}

		_, err := camService.DescribeRoleById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("cam role still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamRoleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cam role %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cam role id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, err := camService.DescribeRoleById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccCamRole_basic = `
resource "tencentcloud_cam_role" "role_basic" {
	name          = "cam-role-test1"
	document      = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/100009461222:uin/100009461222\"]}}]}"
	description   = "test"
	console_login = true
}
`

const testAccCamRole_update = `
resource "tencentcloud_cam_role" "role_basic" {
	name          = "cam-role-test1"
	document      = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/100009461222:uin/100009461222\"]}},{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/100009461222:uin/100009461222\"]}}]}"
	console_login = false
}
`
