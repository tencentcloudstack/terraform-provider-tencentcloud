package cdc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdcDedicatedClustersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCdcDedicatedClustersDataSource,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.tencentcloud_cdc_dedicated_clusters.example", "id"),
			),
		}},
	})
}

const testAccCdcDedicatedClustersDataSource = `
data "tencentcloud_cdc_dedicated_clusters" "example" {}
`
