Provides a resource to create a vpc ccn_instances_reject_attach, you can use this resource to approve cross-region attachment.

Example Usage

```hcl
resource "tencentcloud_ccn_instances_reject_attach" "ccn_instances_reject_attach" {
  ccn_id = "ccn-39lqkygf"
  instances {
    instance_id     = "vpc-j9yhbzpn"
    instance_region = "ap-guangzhou"
    instance_type   = "VPC"
  }
}
```