Provides a resource to create a tdmqRocketmq namespace

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example_namespace"
  remark         = "remark."
}
```
Import

tdmqRocketmq namespace can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_namespace.namespace namespace_id
```