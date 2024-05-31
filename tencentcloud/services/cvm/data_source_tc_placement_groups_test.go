package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudPlacementGroupsDataSource -v
func TestAccTencentCloudPlacementGroupsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementGroupDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementGroupExists("tencentcloud_placement_group.example"),
					resource.TestCheckResourceAttr("data.tencentcloud_placement_groups.placement_group", "placement_group_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_placement_groups.placement_group", "placement_group_list.0.placement_group_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_placement_groups.placement_group", "placement_group_list.0.name", "tf_example"),
					resource.TestCheckResourceAttr("data.tencentcloud_placement_groups.placement_group", "placement_group_list.0.type", "HOST"),
				),
			},
		},
	})
}

const testAccPlacementGroupDataSource = `
resource "tencentcloud_placement_group" "example" {
  name = "tf_example"
  type = "HOST"
}

data "tencentcloud_placement_groups" "placement_group" {
  name               = tencentcloud_placement_group.example.name
  placement_group_id = tencentcloud_placement_group.example.id
}
`
