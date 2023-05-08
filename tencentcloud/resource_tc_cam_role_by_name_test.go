package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCamRoleByNameResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCamRoleByNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRoleByName_basic(ownerUin),
				Check: resource.ComposeTestCheckFunc(
					testAccCamRoleByNameExists("tencentcloud_cam_role_by_name.role_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role_by_name.role_basic", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role_by_name.role_basic", "document"),
				),
			}, {
				Config: testAccCamRoleByName_update(ownerUin),
				Check: resource.ComposeTestCheckFunc(
					testAccCamRoleByNameExists("tencentcloud_cam_role_by_name.role_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role_by_name.role_basic", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role_by_name.role_basic", "document"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_role_by_name.role_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCamRoleByNameDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_role_by_name" {
			continue
		}

		params := make(map[string]interface{})
		params["name"] = rs.Primary.ID
		instances, err := camService.DescribeRolesByFilter(ctx, params)
		if err == nil && len(instances) > 0 {
			return fmt.Errorf("[CHECK][CAM role][Destroy] check: CAM role still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCamRoleByNameExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CAM role][Exists] check: CAM role %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM role][Exists] check: CAM role id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		params := make(map[string]interface{})
		params["name"] = rs.Primary.ID
		instances, err := camService.DescribeRolesByFilter(ctx, params)
		if err != nil {
			return err
		}
		if len(instances) == 0 {
			return fmt.Errorf("[CHECK][CAM role][Exists] check: CAM role %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCamRoleByName_basic(uin string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cam_role_by_name" "role_basic" {
	name          = "cam_role_name_as_identifier_test"
	document      = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/%s:uin/%s\"]}}]}"
	description   = "test"
	console_login = true
}`, uin, uin)
}

func testAccCamRoleByName_update(uin string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cam_role_by_name" "role_basic" {
  name          = "cam_role_name_as_identifier_test"
  document      = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/%s:uin/%s\"]}},{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/%s:uin/%s\"]}}]}"
  console_login = false
}`, uin, uin, uin, uin)
}
