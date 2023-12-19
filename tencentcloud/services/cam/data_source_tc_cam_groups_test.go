package cam_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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

const testAccCamGroupsDataSource_basic = tcacctest.DefaultCamVariables + `
data "tencentcloud_cam_groups" "groups" {
  name = var.cam_group_basic
}
`
