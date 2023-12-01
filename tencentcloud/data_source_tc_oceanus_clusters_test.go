package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusClustersDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusClustersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusClustersDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_clusters.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_clusters.example", "cluster_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_clusters.example", "order_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_clusters.example", "work_space_id"),
				),
			},
		},
	})
}

const testAccOceanusClustersDataSource = `
data "tencentcloud_oceanus_clusters" "example" {
  cluster_ids = ["cluster-5c42n3a5"]
  order_type  = 1
  filters {
    name   = "name"
    values = ["tf_example"]
  }
  work_space_id = "space-2idq8wbr"
}
`
