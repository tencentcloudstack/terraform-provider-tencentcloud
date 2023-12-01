package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudChdfsAccessGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsAccessGroupsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_chdfs_access_groups.access_groups")),
			},
		},
	})
}

const testAccChdfsAccessGroupsDataSource = `

data "tencentcloud_chdfs_access_groups" "access_groups" {
  vpc_id = "vpc-pewdpc0d"
}

`
