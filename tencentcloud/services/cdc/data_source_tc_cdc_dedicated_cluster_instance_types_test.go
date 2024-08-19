package cdc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCdcDedicatedClusterInstanceTypesDataSource_basic -v
func TestAccTencentCloudNeedFixCdcDedicatedClusterInstanceTypesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdcDedicatedClusterInstanceTypesDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdc_dedicated_cluster_instance_types.types", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cdc_dedicated_cluster_instance_types.types", "dedicated_cluster_id"),
				),
			},
		},
	})
}

const testAccCdcDedicatedClusterInstanceTypesDataSource = `
data "tencentcloud_cdc_dedicated_cluster_instance_types" "types" {
  dedicated_cluster_id = "cluster-262n63e8"
}
`
