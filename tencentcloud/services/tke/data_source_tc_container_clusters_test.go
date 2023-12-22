package tke_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDataSourceContainerClusters(t *testing.T) {
	t.Parallel()
	key := "data.tencentcloud_container_clusters.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceContainerClustersConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(key),
					resource.TestCheckResourceAttrSet(key, "total_count"),
					resource.TestCheckResourceAttrSet(key, "clusters.#"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceContainerClustersConfig_basic = `
data "tencentcloud_container_clusters" "foo" {
}
`
