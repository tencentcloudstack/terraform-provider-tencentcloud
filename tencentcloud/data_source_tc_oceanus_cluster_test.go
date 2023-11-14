package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOceanusClusterDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusClusterDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_cluster.cluster")),
			},
		},
	})
}

const testAccOceanusClusterDataSource = `

data "tencentcloud_oceanus_cluster" "cluster" {
  cluster_ids = 
  order_type = 1
  filters {
		name = "name"
		values = 

  }
  work_space_id = "space-1239"
}

`
