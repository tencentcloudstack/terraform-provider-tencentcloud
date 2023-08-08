package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRocketmqNamespaceDataSource -v
func TestAccTencentCloudTdmqRocketmqNamespaceDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTdmqRocketmqNamespace,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rocketmq_namespace.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdmq_rocketmq_namespace.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdmq_rocketmq_namespace.example", "name_keyword"),
				),
			},
		},
	})
}

const testAccDataSourceTdmqRocketmqNamespace = `
data "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id   = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  name_keyword = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
}

resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example"
  remark         = "remark."
}
`
