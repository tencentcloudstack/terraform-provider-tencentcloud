Provides a resource to create a dnspod record_group

Example Usage

```hcl
resource "tencentcloud_dnspod_record_group" "record_group" {
  domain = "dnspod.cn"
  group_name = "group_demo"
}
```

Import

dnspod record_group can be imported using the domain#groupId, e.g.

```
terraform import tencentcloud_dnspod_record_group.record_group domain#groupId
```