Provides a resource to create a tdmqRocketmq environment_role

Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_role" "example" {
  role_name  = "tf_example_role"
  remark     = "remark."
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example_namespace"
  remark         = "remark."
}

resource "tencentcloud_tdmq_rocketmq_environment_role" "example" {
  environment_name = tencentcloud_tdmq_rocketmq_namespace.example.namespace_name
  role_name        = tencentcloud_tdmq_rocketmq_role.example.role_name
  permissions      = ["produce", "consume"]
  cluster_id       = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
}
```
Import

tdmqRocketmq environment_role can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_environment_role.environment_role environmentRole_id
```