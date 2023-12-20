package cam_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cam"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCamGroupPolicyAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCamGroupPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamGroupPolicyAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamGroupPolicyAttachmentExists("tencentcloud_cam_group_policy_attachment.group_policy_attachment_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_group_policy_attachment.group_policy_attachment_basic", "group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_group_policy_attachment.group_policy_attachment_basic", "policy_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_group_policy_attachment.group_policy_attachment_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamGroupPolicyAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	camService := cam.NewCamService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_group_policy_attachment" {
			continue
		}

		instance, err := camService.DescribeGroupPolicyAttachmentById(ctx, rs.Primary.ID)
		if err == nil && instance != nil {
			return fmt.Errorf("[CHECK][CAM group policy attachment][Destroy] check: CAM group policy attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamGroupPolicyAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CAM group policy attachment][Exists] check: CAM group policy attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM group policy attachment][Exists] check: CAM group policy attachment id is not set")
		}
		camService := cam.NewCamService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := camService.DescribeGroupPolicyAttachmentById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CAM group policy attachment][Exists] check: CAM group policy attachment %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCamGroupPolicyAttachment_basic = tcacctest.DefaultCamVariables + `
data "tencentcloud_cam_groups" "groups" {
  name = var.cam_group_basic
}

data "tencentcloud_cam_policies" "policy" {
  name = var.cam_policy_basic
}

resource "tencentcloud_cam_group_policy_attachment" "group_policy_attachment_basic" {
  group_id  = data.tencentcloud_cam_groups.groups.group_list.0.group_id
  policy_id = data.tencentcloud_cam_policies.policy.policy_list.0.policy_id
}
`
