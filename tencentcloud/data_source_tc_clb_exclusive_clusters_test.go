package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClbExclusiveClustersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbExclusiveClustersDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clb_exclusive_clusters.exclusive_clusters")),
			},
		},
	})
}

const testAccClbExclusiveClustersDataSource = `

data "tencentcloud_clb_exclusive_clusters" "exclusive_clusters" {
  filters {
		name = ""
		values = 

  }
  }

`
