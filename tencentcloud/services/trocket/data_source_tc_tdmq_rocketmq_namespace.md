Use this data source to query detailed information of tdmqRocketmq namespace

Example Usage

```hcl
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
```