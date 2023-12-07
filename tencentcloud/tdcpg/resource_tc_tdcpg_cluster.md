Provides a resource to create a tdcpg cluster.

~> **NOTE:** This resource is still in internal testing. To experience its functions, you need to apply for a whitelist from Tencent Cloud.

Example Usage

```hcl
resource "tencentcloud_tdcpg_cluster" "cluster" {
  zone = "ap-guangzhou-3"
  master_user_password = ""
  cpu = 1
  memory = 1
  vpc_id = "vpc_id"
  subnet_id = "subnet_id"
  pay_mode = "POSTPAID_BY_HOUR"
  cluster_name = "cluster_name"
  db_version = "10.17"
  instance_count = 1
  period = 1
  project_id = 0
}

```
Import

tdcpg cluster can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdcpg_cluster.cluster cluster_id
```