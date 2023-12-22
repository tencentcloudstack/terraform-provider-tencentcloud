package trocket_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRocketmqTopicDataSource -v
func TestAccTencentCloudTdmqRocketmqTopicDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTdmqRocketmqTopic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rocketmq_topic.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdmq_rocketmq_topic.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdmq_rocketmq_topic.example", "namespace_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdmq_rocketmq_topic.example", "filter_name"),
				),
			},
		},
	})
}

const testAccDataSourceTdmqRocketmqTopic = `
data "tencentcloud_tdmq_rocketmq_topic" "example" {
  cluster_id   = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_id = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  filter_name  = tencentcloud_tdmq_rocketmq_topic.example.topic_name
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

resource "tencentcloud_tdmq_rocketmq_topic" "example" {
  topic_name     = "tf_example"
  namespace_name = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  type           = "Normal"
  remark         = "remark."
}
`
