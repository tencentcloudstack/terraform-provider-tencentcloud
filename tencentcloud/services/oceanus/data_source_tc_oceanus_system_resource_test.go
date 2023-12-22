package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusSystemResourceDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusSystemResourceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusSystemResourceDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_system_resource.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_system_resource.example", "resource_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_system_resource.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_system_resource.example", "flink_version"),
				),
			},
		},
	})
}

const testAccOceanusSystemResourceDataSource = `
data "tencentcloud_oceanus_system_resource" "example" {
  resource_ids = ["resource-abd503yt"]
  filters {
    name   = "Name"
    values = ["tf_example"]
  }
  cluster_id    = "cluster-n8yaia0p"
  flink_version = "Flink-1.11"
}
`
