Provides a resource to create an identity center user

Example Usage

```hcl
resource "tencentcloud_identity_center_user" "identity_center_user" {
    zone_id = "z-xxxxxx"
    user_name = "test-user"
    description = "test"
}
```

Import

organization identity_center_user can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_user.identity_center_user ${zoneId}#${userId}
```
