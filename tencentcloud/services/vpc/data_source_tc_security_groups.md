Use this data source to query detailed information of security groups.

Example Usage

```hcl
data "tencentcloud_security_groups" "sglab" {
  security_group_id = tencentcloud_security_group.sglab.id
}
```