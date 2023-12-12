Provides a resource to create a vpc ccn_instances_reset_attach, you can use this resource to reset cross-region attachment.

Example Usage

```hcl
resource "tencentcloud_ccn_instances_reset_attach" "ccn_instances_reset_attach" {
  ccn_id = "ccn-39lqkygf"
  ccn_uin = "100022975249"
  instances {
    instance_id = "vpc-j9yhbzpn"
    instance_region = "ap-guangzhou"
    instance_type = "VPC"
  }
}
```