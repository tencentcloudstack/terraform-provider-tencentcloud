package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCynosdbClusterInstanceGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterInstanceGroupsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_cluster_instance_groups.cluster_instance_groups"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_cluster_instance_groups.cluster_instance_groups", "instance_grp_info_list.#"),
				),
			},
		},
	})
}

const testAccCynosdbClusterInstanceGroupsDataSource = CommonCynosdb + `

data "tencentcloud_cynosdb_cluster_instance_groups" "cluster_instance_groups" {
	cluster_id = var.cynosdb_cluster_id
}

`
