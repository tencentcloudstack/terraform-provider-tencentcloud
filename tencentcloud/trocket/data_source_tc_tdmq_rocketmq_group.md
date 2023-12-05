Use this data source to query detailed information of tdmqRocketmq group

Example Usage

```hcl
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
```