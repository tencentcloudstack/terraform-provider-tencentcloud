package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCamGroupMembershipsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamGroupMembershipDestroy,
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

const testAccCamGroupMembershipsDataSource_basic = defaultCamVariables + `
data "tencentcloud_cam_groups" "groups" {
  name = var.cam_group_basic
}

data "tencentcloud_cam_users" "users" {
  name = var.cam_user_basic
}

resource "tencentcloud_cam_group_membership" "membership" {
  group_id = data.tencentcloud_cam_groups.groups.group_list.0.group_id
  user_names = [data.tencentcloud_cam_users.users.user_list.0.user_id]
}

data "tencentcloud_cam_group_memberships" "memberships" {
  group_id = tencentcloud_cam_group_membership.membership.id
}
`
