package trocket_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRocketmqGroupDataSource -v
func TestAccTencentCloudTdmqRocketmqGroupDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTdmqRocketmqGroup,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rocketmq_group.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdmq_rocketmq_group.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdmq_rocketmq_group.example", "namespace_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdmq_rocketmq_group.example", "filter_group"),
				),
			},
		},
	})
}

const testAccDataSourceTdmqRocketmqGroup = `
data "tencentcloud_tdmq_rocketmq_group" "example" {
  cluster_id   = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_id = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  filter_group = tencentcloud_tdmq_rocketmq_group.example.group_name
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

resource "tencentcloud_tdmq_rocketmq_group" "example" {
  group_name       = "tf_example"
  namespace        = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  read_enable      = true
  broadcast_enable = true
  cluster_id       = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  remark           = "remark."
}
`
