
package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTdcpgClustersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTdcpgClusters,
				Check: resource.ComposeTestCheckFunc(
				  testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdcpg_clusters.clusters"),
				),
			},
		},
	})
}

const testAccDataSourceTdcpgClusters = `

data "tencentcloud_tdcpg_clusters" "clusters" {
  cluster_id = ""
  cluster_name = ""
  status = ""
  pay_mode = ""
  project_id = ""
  }

`
