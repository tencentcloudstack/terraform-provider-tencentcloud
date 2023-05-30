package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClbClusterResourcesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbClusterResourcesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clb_cluster_resources.cluster_resources")),
			},
		},
	})
}

const testAccClbClusterResourcesDataSource = `

data "tencentcloud_clb_cluster_resources" "cluster_resources" {
  filters {
    name = "idle"
    values = ["True"]
  }
}

`
