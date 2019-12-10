package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCamRolePolicyAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamRolePolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRolePolicyAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamRolePolicyAttachmentExists("tencentcloud_cam_role_policy_attachment.role_policy_attachment_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role_policy_attachment.role_policy_attachment_basic", "role_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role_policy_attachment.role_policy_attachment_basic", "policy_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_role_policy_attachment.role_policy_attachment_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamRolePolicyAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_role_policy_attachment" {
			continue
		}

		_, err := camService.DescribeRolePolicyAttachmentById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("CAM role policy attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamRolePolicyAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("CAM role policy attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("CAM role policy attachment id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, err := camService.DescribeRolePolicyAttachmentById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

//need to add policy resource definition
const testAccCamRolePolicyAttachment_basic = `
resource "tencentcloud_cam_role" "role" {
  name          = "cam-role-test"
  document      = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/100009461222:uin/100009461222\"]}}]}"
  description   = "test"
  console_login = true
}

resource "tencentcloud_cam_policy" "policy" {
  name        = "cam-policy-test2"
  document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
  description = "test"
}
  
resource "tencentcloud_cam_role_policy_attachment" "role_policy_attachment_basic" {
  role_id   = "${tencentcloud_cam_role.role.id}"
  policy_id = "${tencentcloud_cam_policy.policy.id}"
}
`
