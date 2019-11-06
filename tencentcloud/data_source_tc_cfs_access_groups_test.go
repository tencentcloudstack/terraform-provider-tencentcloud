package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudCfsAccessGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCfsAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsAccessGroupsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCfsAccessGroupExists("tencentcloud_cfs_access_group.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_access_groups.access_groups", "access_group_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_access_groups.access_groups", "access_group_list.0.name", "test_cfs_access_group"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_access_groups.access_groups", "access_group_list.0.description", "test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_access_groups.access_groups", "access_group_list.0.access_group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_access_groups.access_groups", "access_group_list.0.create_time"),
				),
			},
		},
	})
}

const testAccCfsAccessGroupsDataSource = `
resource "tencentcloud_cfs_access_group" "foo" {
  name = "test_cfs_access_group"
  description = "test"
}

data "tencentcloud_cfs_access_groups" "access_groups" {
  access_group_id = "${tencentcloud_cfs_access_group.foo.id}"
  name = "${tencentcloud_cfs_access_group.foo.name}"
}
`
