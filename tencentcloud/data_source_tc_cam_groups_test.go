package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCamGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamGroupsDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_groups.groups", "group_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_groups.groups", "group_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_groups.groups", "group_list.0.remark"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_groups.groups", "group_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_groups.groups", "group_list.0.group_id"),
				),
			},
		},
	})
}

const testAccCamGroupsDataSource_basic = defaultCamVariables + `
data "tencentcloud_cam_groups" "groups" {
  name = var.cam_group_basic
}
`
