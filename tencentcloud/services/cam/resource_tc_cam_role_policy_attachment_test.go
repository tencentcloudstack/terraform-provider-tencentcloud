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

// go test -i; go test -test.run TestAccTencentCloudCamRolePolicyAttachmentResource_basic -v
func TestAccTencentCloudCamRolePolicyAttachmentResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCamRolePolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRolePolicyAttachment_basic(),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	camService := cam.NewCamService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_role_policy_attachment" {
			continue
		}

		instance, err := camService.DescribeRolePolicyAttachmentById(ctx, rs.Primary.ID)
		if err == nil && instance != nil {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Desctroy] check: CAM role policy attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamRolePolicyAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Exist] check: CAM role policy attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Exist] check: CAM role policy attachment id is not set")
		}
		camService := cam.NewCamService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := camService.DescribeRolePolicyAttachmentById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Exist] check: CAM role policy attachment %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

// need to add policy resource definition
func testAccCamRolePolicyAttachment_basic() string {
	return tcacctest.DefaultCamVariables + `
data "tencentcloud_cam_policies" "policy" {
  name        = var.cam_policy_basic
}

data "tencentcloud_cam_roles" "roles" {
  name        = var.cam_role_basic
}

resource "tencentcloud_cam_role_policy_attachment" "role_policy_attachment_basic" {
  role_id   = data.tencentcloud_cam_roles.roles.role_list.0.role_id
  policy_id = data.tencentcloud_cam_policies.policy.policy_list.0.policy_id
}

`

}
