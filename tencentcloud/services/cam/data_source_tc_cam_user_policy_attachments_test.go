package cam_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamUserPolicyAttachmentsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCamUserPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUserPolicyAttachmentsDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_cam_user_policy_attachments.user_policy_attachments", "user_policy_attachment_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_user_policy_attachments.user_policy_attachments", "user_policy_attachment_list.0.user_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_user_policy_attachments.user_policy_attachments", "user_policy_attachment_list.0.policy_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_user_policy_attachments.user_policy_attachments", "user_policy_attachment_list.0.policy_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_user_policy_attachments.user_policy_attachments", "user_policy_attachment_list.0.create_mode"),
				),
			},
		},
	})
}

const testAccCamUserPolicyAttachmentsDataSource_basic = tcacctest.DefaultCamVariables + `

data "tencentcloud_cam_user_policy_attachments" "user_policy_attachments" {
  user_name = var.cam_user_basic
}
`
