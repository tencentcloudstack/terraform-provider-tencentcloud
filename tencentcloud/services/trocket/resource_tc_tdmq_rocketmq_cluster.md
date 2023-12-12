Provides a resource to create a tdmqRocketmq cluster

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}
```
Import

tdmqRocketmq cluster can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_cluster.cluster cluster_id
```