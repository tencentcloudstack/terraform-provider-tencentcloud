package cdc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCdcDedicatedClusterOrdersDataSource_basic -v
func TestAccTencentCloudNeedFixCdcDedicatedClusterOrdersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdcDedicatedClusterOrdersDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdc_dedicated_cluster_orders.orders", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdc_dedicated_cluster_orders.orders", "dedicated_cluster_ids.#"),
				),
			},
		},
	})
}

const testAccCdcDedicatedClusterOrdersDataSource = `
data "tencentcloud_cdc_dedicated_cluster_orders" "orders1" {
  dedicated_cluster_ids = ["cluster-262n63e8"]
}
`
