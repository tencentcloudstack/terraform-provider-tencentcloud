Provides a resource to create an identity center group

Example Usage

```hcl
resource "tencentcloud_identity_center_group" "identity_center_group" {
    zone_id = "z-xxxxxx"
    group_name = "test-group"
    description = "test"
}
```

Import

tencentcloud_identity_center_group can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_group.identity_center_group ${zoneId}#${groupId}
```
