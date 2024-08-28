Provides a resource to create an identity center user group attachment

Example Usage

```hcl
resource "tencentcloud_identity_center_user_group_attachment" "identity_center_user_group_attachment" {
    zone_id = "z-xxxxxx"
    user_id = "u-xxxxxx"
    group_id = "g-xxxxxx"
}
```

Import

organization identity_center_user_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_user_group_attachment.identity_center_user_group_attachment ${zoneId}#${groupId}#${userId}
```
