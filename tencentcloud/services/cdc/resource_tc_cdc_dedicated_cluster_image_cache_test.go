package cdc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdcDedicatedClusterImageCacheResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdcDedicatedCluster,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster_image_cache.cdc_dedicated_cluster_image_cache", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster_image_cache.cdc_dedicated_cluster_image_cache", "dedicated_cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster_image_cache.cdc_dedicated_cluster_image_cache", "image_id"),
				),
			},
		},
	})
}

const testAccCdcDedicatedClusterImageCache = `
resource "tencentcloud_cdc_dedicated_cluster_image_cache" "cdc_dedicated_cluster_image_cache" {
  dedicated_cluster_id = ""
  image_id = ""
}
`
