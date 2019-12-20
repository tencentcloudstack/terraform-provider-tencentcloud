package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCamGroupPolicyAttachmentsDataSource_basic(t *testing.T) {
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

const testAccCamGroupPolicyAttachmentsDataSource_basic = `
resource "tencentcloud_cam_group" "group" {
  name   = "cam-group-policy-test"
  remark = "test"
}

resource "tencentcloud_cam_policy" "policy" {
  name        = "cam-policy-test8"
  document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
  description = "test"
}

resource "tencentcloud_cam_group_policy_attachment" "group_policy_attachment" {
  group_id  = "${tencentcloud_cam_group.group.id}"
  policy_id = "${tencentcloud_cam_policy.policy.id}"
}
  
data "tencentcloud_cam_group_policy_attachments" "group_policy_attachments" {
  group_id = "${tencentcloud_cam_group_policy_attachment.group_policy_attachment.group_id}"
}
`
