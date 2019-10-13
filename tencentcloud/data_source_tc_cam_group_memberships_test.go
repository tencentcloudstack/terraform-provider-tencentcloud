package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudCamGroupMembershipsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamGroupMembershipsDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCamGroupExists("tencentcloud_cam_group_membership.membership"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_group_memberships.memberships", "membership_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_group_memberships.memberships", "membership_list.0.group_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_group_memberships.memberships", "membership_list.0.user_ids.#", "1"),
				),
			},
		},
	})
}

const testAccCamGroupMembershipsDataSource_basic = `
resource "tencentcloud_cam_group" "group_basic" {
	name   = "cam-group-membership-test"
	remark = "test"
}

resource "tencentcloud_cam_user" "user_basic" {
	name                = "cam-user-testj"
	remark              = "test"
	console_login       = true
	use_api             = true
	need_reset_password = true
	password            = "Gail@1234"
	phone_num           = "13631555963"
	country_code        = "86"
	email               = "1234@qq.com"
}

resource "tencentcloud_cam_group_membership" "membership" {
	group_id = "${tencentcloud_cam_group.group_basic.id}"
	user_ids = ["${tencentcloud_cam_user.user_basic.id}"]
}

data "tencentcloud_cam_group_memberships" "memberships" {
	group_id = "${tencentcloud_cam_group_membership.membership.id}"
}
`
