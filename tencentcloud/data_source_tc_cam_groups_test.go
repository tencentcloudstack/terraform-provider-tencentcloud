package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudCamGroupsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamGroupsDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCamGroupExists("tencentcloud_cam_group.group"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_groups.groups", "group_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_groups.groups", "group_list.0.name", "cam-group-test3"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_groups.groups", "group_list.0.remark", "test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_groups.groups", "group_list.0.create_time"),
				),
			},
		},
	})
}

const testAccCamGroupsDataSource_basic = `
resource "tencentcloud_cam_group" "group" {
  name   = "cam-group-test3"
  remark = "test"
}
  
data "tencentcloud_cam_groups" "groups" {
  group_id = "${tencentcloud_cam_group.group.id}"
}
`
