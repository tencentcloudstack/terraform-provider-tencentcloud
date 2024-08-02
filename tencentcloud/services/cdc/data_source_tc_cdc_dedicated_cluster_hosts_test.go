package cdc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCdcDedicatedClusterHostsDataSource_basic -v
func TestAccTencentCloudNeedFixCdcDedicatedClusterHostsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdcDedicatedClusterHostsDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdc_dedicated_cluster_hosts.hosts", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdc_dedicated_cluster_hosts.hosts", "dedicated_cluster_id"),
				),
			},
		},
	})
}

const testAccCdcDedicatedClusterHostsDataSource = `
data "tencentcloud_cdc_dedicated_cluster_hosts" "hosts" {
  dedicated_cluster_id = "cluster-262n63e8"
}
`
