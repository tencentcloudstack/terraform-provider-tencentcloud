package clb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClbExclusiveClustersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbExclusiveClustersDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_clb_exclusive_clusters.exclusive_clusters")),
			},
		},
	})
}

const testAccClbExclusiveClustersDataSource = `

data "tencentcloud_clb_exclusive_clusters" "exclusive_clusters" {
  filters {
    name = "zone"
    values = ["ap-guangzhou-1"]
  }
}

`
