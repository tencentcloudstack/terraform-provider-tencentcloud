package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCamRolePolicyAttachmentsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamRolePolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRolePolicyAttachmentsDataSource_basic(ownerUin),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCamRolePolicyAttachmentExists("tencentcloud_cam_role_policy_attachment.role_policy_attachment"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_role_policy_attachments.role_policy_attachments", "role_policy_attachment_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_role_policy_attachments.role_policy_attachments", "role_policy_attachment_list.0.role_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_role_policy_attachments.role_policy_attachments", "role_policy_attachment_list.0.policy_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_role_policy_attachments.role_policy_attachments", "role_policy_attachment_list.0.policy_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_role_policy_attachments.role_policy_attachments", "role_policy_attachment_list.0.create_mode"),
				),
			},
		},
	})
}

func testAccCamRolePolicyAttachmentsDataSource_basic(uin string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cam_role" "role" {
  name          = "cam-role-test"
  document      = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/%s:uin/%s\"]}}]}"
  description   = "test"
  console_login = true
}

resource "tencentcloud_cam_policy" "policy" {
  name        = "cam-policy-test6"
  document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
  description = "test"
}

resource "tencentcloud_cam_role_policy_attachment" "role_policy_attachment" {
  role_id   = tencentcloud_cam_role.role.id
  policy_id = tencentcloud_cam_policy.policy.id
}
  
data "tencentcloud_cam_role_policy_attachments" "role_policy_attachments" {
  role_id = tencentcloud_cam_role_policy_attachment.role_policy_attachment.role_id
}`, uin, uin)
}
