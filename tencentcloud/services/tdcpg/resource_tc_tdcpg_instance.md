Provides a resource to create a tdcpg instance.

~> **NOTE:** This resource is still in internal testing. To experience its functions, you need to apply for a whitelist from Tencent Cloud.

Example Usage

```hcl
resource "tencentcloud_tdcpg_instance" "instance1" {
  cluster_id = "cluster_id"
  cpu = 1
  memory = 1
  instance_name = "instance_name"
}

resource "tencentcloud_tdcpg_instance" "instance2" {
  cluster_id = "cluster_id"
  cpu = 1
  memory = 2
  instance_name = "instance_name"
  operation_timing = "IMMEDIATE"
}

```
Import

tdcpg instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdcpg_instance.instance cluster_id#instance_id
```