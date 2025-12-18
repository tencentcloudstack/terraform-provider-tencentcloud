Provides a resource to create a BH resource

~> **NOTE:** Currently, executing the `terraform destroy` command to delete this resource is not supported. If you need to destroy it, please contact Tencent Cloud BH through a ticket.

Example Usage

```hcl
resource "tencentcloud_bh_resource" "example" {
  deploy_region    = "ap-guangzhou"
  vpc_id           = "vpc-q1of50wz"
  subnet_id        = "subnet-7uhvm46o"
  resource_edition = "standard"
  resource_node    = 20
  time_unit        = "m"
  time_span        = "1"
  pay_mode         = 1
  auto_renew_flag  = 1
  deploy_zone      = "ap-guangzhou-6"
  cidr_block       = "192.168.11.0/24"
  vpc_cidr_block   = "192.168.0.0/16"
  intranet_access  = 1
  external_access  = 1
}
```
