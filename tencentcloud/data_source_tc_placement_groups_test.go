package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPlacementGroupsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementGroupDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementGroupExists("tencentcloud_placement_group.placement"),
					resource.TestCheckResourceAttr("data.tencentcloud_placement_groups.data_placement", "placement_group_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_placement_groups.data_placement", "placement_group_list.0.placement_group_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_placement_groups.data_placement", "placement_group_list.0.name", "tf-test-placement"),
					resource.TestCheckResourceAttr("data.tencentcloud_placement_groups.data_placement", "placement_group_list.0.type", "HOST"),
				),
			},
		},
	})
}

const testAccPlacementGroupDataSource = `
resource "tencentcloud_placement_group" "placement" {
  name = "tf-test-placement"
  type = "HOST"
}

data "tencentcloud_placement_groups" "data_placement" {
  placement_group_id = tencentcloud_placement_group.placement.id
}
`
