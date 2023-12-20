package cfs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCfsAccessGroupsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCfsAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsAccessGroupsDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_access_groups.access_groups", "access_group_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_access_groups.access_groups", "access_group_list.0.name", "keep_access_group"),
					resource.TestCheckResourceAttr("data.tencentcloud_cfs_access_groups.access_groups", "access_group_list.0.description", "test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_access_groups.access_groups", "access_group_list.0.access_group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cfs_access_groups.access_groups", "access_group_list.0.create_time"),
				),
			},
		},
	})
}

const DefaultCfsAccessGroup = `
data "tencentcloud_cfs_access_groups" "access_groups" {
  name = "keep_access_group"
}

locals {
  cfs_access_group_id = data.tencentcloud_cfs_access_groups.access_groups.access_group_list.0.access_group_id
}
`

const testAccCfsAccessGroupsDataSource = `
data "tencentcloud_cfs_access_groups" "access_groups" {
  name = "keep_access_group"
}
`
