Use this data source to query detailed information of tdmqRocketmq topic

Example Usage

```hcl
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
```