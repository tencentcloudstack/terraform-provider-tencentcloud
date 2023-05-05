package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamRolePolicyAttachmentsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamRolePolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRolePolicyAttachmentsDataSource_basic,
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

const testAccCamRolePolicyAttachmentsDataSource_basic = defaultCamVariables + `
data "tencentcloud_cam_policies" "policy" {
  name        = var.cam_policy_basic
}

data "tencentcloud_cam_roles" "roles" {
  name        = var.cam_role_basic
}

resource "tencentcloud_cam_role_policy_attachment" "role_policy_attachment" {
  role_id   = data.tencentcloud_cam_roles.roles.role_list.0.role_id
  policy_id = data.tencentcloud_cam_policies.policy.policy_list.0.policy_id
}

data "tencentcloud_cam_role_policy_attachments" "role_policy_attachments" {
  role_id = tencentcloud_cam_role_policy_attachment.role_policy_attachment.role_id
}`
