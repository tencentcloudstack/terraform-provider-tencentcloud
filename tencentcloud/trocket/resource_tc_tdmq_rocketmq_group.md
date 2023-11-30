Provides a resource to create a tdmqRocketmq group

Example Usage

```hcl
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
  cluster_id       = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace        = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  read_enable      = true
  broadcast_enable = true
  remark           = "remark."
}
```
Import

tdmqRocketmq group can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_group.group group_id
```