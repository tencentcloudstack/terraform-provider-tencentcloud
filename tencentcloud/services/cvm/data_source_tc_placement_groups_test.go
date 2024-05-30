package cvm_test

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudPlacementGroupsDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementGroupsDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_placement_groups.data_placement"), resource.TestCheckResourceAttrSet("data.tencentcloud_placement_groups.data_placement", "placement_group_list.0.placement_group_id"), resource.TestCheckResourceAttr("data.tencentcloud_placement_groups.data_placement", "placement_group_list.0.type", "HOST"), resource.TestCheckResourceAttr("data.tencentcloud_placement_groups.data_placement", "placement_group_list.0.name", "tf-test-placement"), resource.TestCheckResourceAttr("data.tencentcloud_placement_groups.data_placement", "placement_group_list.#", "1")),
			},
		},
	})
}

const testAccPlacementGroupsDataSource_BasicCreate = `

data "tencentcloud_placement_groups" "data_placement" {
    placement_group_id = tencentcloud_placement_group.placement.id
}
resource "tencentcloud_placement_group" "placement" {
    name = "tf-test-placement"
    type = "HOST"
}

`
