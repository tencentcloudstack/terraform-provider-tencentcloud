package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbClusterDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_cluster.cluster")),
			},
		},
	})
}

const testAccCynosdbClusterDataSource = `

data "tencentcloud_cynosdb_cluster" "cluster" {
  cluster_id = "xxx"
  database = "test"
  table = "1"
  table_type = "all"
  }

`
