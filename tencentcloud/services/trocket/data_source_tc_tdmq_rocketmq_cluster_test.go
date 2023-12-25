package trocket_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRocketmqClusterDataSource -v
func TestAccTencentCloudTdmqRocketmqClusterDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRocketmqCluster,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rocketmq_cluster.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdmq_rocketmq_cluster.example", "name_keyword"),
				),
			},
		},
	})
}

const testAccDataSourceRocketmqCluster = `
data "tencentcloud_tdmq_rocketmq_cluster" "example" {
  name_keyword = tencentcloud_tdmq_rocketmq_cluster.example.cluster_name
}

resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}
`
