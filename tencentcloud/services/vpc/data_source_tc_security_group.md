Use this data source to query detailed information of security group.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_security_groups.

Example Usage

```hcl
data "tencentcloud_security_group" "sglab" {
  security_group_id = tencentcloud_security_group.sglab.id
}
```