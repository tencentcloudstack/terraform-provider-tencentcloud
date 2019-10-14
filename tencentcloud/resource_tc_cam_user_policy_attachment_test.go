package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCamUserPolicyAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamUserPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUserPolicyAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamUserPolicyAttachmentExists("tencentcloud_cam_user_policy_attachment.user_policy_attachment_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user_policy_attachment.user_policy_attachment_basic", "user_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user_policy_attachment.user_policy_attachment_basic", "policy_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_user_policy_attachment.user_policy_attachment_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamUserPolicyAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_user_policy_attachment" {
			continue
		}

		_, err := camService.DescribeUserPolicyAttachmentById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("cam user policy attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamUserPolicyAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cam user policy attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cam user policy attachment id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, err := camService.DescribeUserPolicyAttachmentById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

//need to add policy resource definition
const testAccCamUserPolicyAttachment_basic = `
resource "tencentcloud_cam_user" "user" {
  name                = "cam-user-testt"
  remark              = "test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "13631555963"
  country_code        = "86"
  email               = "1234@qq.com"
}

resource "tencentcloud_cam_policy" "policy" {
  name        = "cam-policy-test3"
  document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
  description = "test"
}

resource "tencentcloud_cam_user_policy_attachment" "user_policy_attachment_basic" {
  user_id   = "${tencentcloud_cam_user.user.id}"
  policy_id = "${tencentcloud_cam_policy.policy.id}"
}
`
