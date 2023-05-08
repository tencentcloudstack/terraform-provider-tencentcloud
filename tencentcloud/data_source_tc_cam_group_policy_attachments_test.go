package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamGroupPolicyAttachmentsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamGroupPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamGroupPolicyAttachmentsDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCamGroupPolicyAttachmentExists("tencentcloud_cam_group_policy_attachment.group_policy_attachment"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_group_policy_attachments.group_policy_attachments", "group_policy_attachment_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_group_policy_attachments.group_policy_attachments", "group_policy_attachment_list.0.group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_group_policy_attachments.group_policy_attachments", "group_policy_attachment_list.0.policy_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_group_policy_attachments.group_policy_attachments", "group_policy_attachment_list.0.policy_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_group_policy_attachments.group_policy_attachments", "group_policy_attachment_list.0.create_mode"),
				),
			},
		},
	})
}

const testAccCamGroupPolicyAttachmentsDataSource_basic = defaultCamVariables + `
data "tencentcloud_cam_groups" "groups" {
  name = var.cam_group_basic
}

data "tencentcloud_cam_policies" "policy" {
  name        = var.cam_policy_basic
}


resource "tencentcloud_cam_group_policy_attachment" "group_policy_attachment" {
  group_id  = data.tencentcloud_cam_groups.groups.group_list.0.group_id
  policy_id = data.tencentcloud_cam_policies.policy.policy_list.0.policy_id
}

data "tencentcloud_cam_group_policy_attachments" "group_policy_attachments" {
  group_id = tencentcloud_cam_group_policy_attachment.group_policy_attachment.group_id
}
`
