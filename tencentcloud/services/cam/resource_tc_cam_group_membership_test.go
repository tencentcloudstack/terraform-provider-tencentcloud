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

func TestAccTencentCloudCamGroupMembership_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCamGroupMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamGroupMembership_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamGroupMembershipExists("tencentcloud_cam_group_membership.group_membership_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_group_membership.group_membership_basic", "group_id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_group_membership.group_membership_basic", "user_names.#", "1"),
				),
			}, {
				Config: testAccCamGroupMembership_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamGroupMembershipExists("tencentcloud_cam_group_membership.group_membership_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_group_membership.group_membership_basic", "group_id"),
					resource.TestCheckResourceAttr("tencentcloud_cam_group_membership.group_membership_basic", "user_names.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_group_membership.group_membership_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamGroupMembershipDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	camService := cam.NewCamService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_group_membership" {
			continue
		}

		instance, err := camService.DescribeGroupMembershipById(ctx, rs.Primary.ID)
		if err == nil && len(instance) > 0 {
			return fmt.Errorf("[CHECK][CAM group membership][Destroy] check: CAM group membership still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamGroupMembershipExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CAM group membership][Exists] check: CAM group membership %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM group membership][Exists] check: CAM group membership id is not set")
		}
		camService := cam.NewCamService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := camService.DescribeGroupMembershipById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if len(instance) == 0 {
			return fmt.Errorf("[CHECK][CAM group membership][Exists] check: CAM group membership %s is not exists", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCamGroupMembership_basic = tcacctest.DefaultCamVariables + `
data "tencentcloud_cam_groups" "groups" {
  name   = var.cam_group_basic
}

resource "tencentcloud_cam_user" "foo" {
  name                = "cam-user-test22"
  remark              = "test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "12345678910"
  country_code        = "86"
  email               = "1234@qq.com"
  force_delete        = true
}

resource "tencentcloud_cam_user" "user_basic" {
  name                = "cam-user-test33"
  remark              = "test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "12345678910"
  country_code        = "86"
  email               = "1234@qq.com"
  force_delete        = true
}

resource "tencentcloud_cam_group_membership" "group_membership_basic" {
  group_id = data.tencentcloud_cam_groups.groups.group_list.0.group_id
  user_names = [tencentcloud_cam_user.foo.id]
}
`

const testAccCamGroupMembership_update = tcacctest.DefaultCamVariables + `
data "tencentcloud_cam_groups" "groups" {
  name   = var.cam_group_basic
}

resource "tencentcloud_cam_user" "foo" {
  name                = "cam-user-test22"
  remark              = "test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "12345678910"
  country_code        = "86"
  email               = "1234@qq.com"
  force_delete        = true
}

resource "tencentcloud_cam_user" "user_basic" {
  name                = "cam-user-test33"
  remark              = "test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "12345678910"
  country_code        = "86"
  email               = "1234@qq.com"
  force_delete        = true
}

resource "tencentcloud_cam_group_membership" "group_membership_basic" {
  group_id = data.tencentcloud_cam_groups.groups.group_list.0.group_id
  user_names = [tencentcloud_cam_user.user_basic.id]
}
`
