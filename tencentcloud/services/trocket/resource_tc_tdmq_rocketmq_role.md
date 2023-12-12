Provides a resource to create a tdmqRocketmq role

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_role" "example" {
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  role_name  = "tf_example"
  remark     = "remark."
}
```
Import

tdmqRocketmq role can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_role.role role_id
```