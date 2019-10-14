package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudCamUserPolicyAttachmentsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamUserPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUserPolicyAttachmentsDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCamUserPolicyAttachmentExists("tencentcloud_cam_user_policy_attachment.user_policy_attachment"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_user_policy_attachments.user_policy_attachments", "user_policy_attachment_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_user_policy_attachments.user_policy_attachments", "user_policy_attachment_list.0.user_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_user_policy_attachments.user_policy_attachments", "user_policy_attachment_list.0.policy_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_user_policy_attachments.user_policy_attachments", "user_policy_attachment_list.0.policy_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_user_policy_attachments.user_policy_attachments", "user_policy_attachment_list.0.create_mode"),
				),
			},
		},
	})
}

const testAccCamUserPolicyAttachmentsDataSource_basic = `
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
  name        = "cam-policy-test7"
  document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
  description = "test"
}

resource "tencentcloud_cam_user_policy_attachment" "user_policy_attachment" {
  user_id   = "${tencentcloud_cam_user.user.id}"
  policy_id = "${tencentcloud_cam_policy.policy.id}"
}
  
data "tencentcloud_cam_user_policy_attachments" "user_policy_attachments" {
  user_id = "${tencentcloud_cam_user_policy_attachment.user_policy_attachment.user_id}"
}
`
