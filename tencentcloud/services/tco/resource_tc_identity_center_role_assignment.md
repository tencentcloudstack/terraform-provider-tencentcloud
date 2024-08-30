Provides a resource to create a organization identity_center_role_assignment

Example Usage

```hcl
resource "tencentcloud_identity_center_role_assignment" "identity_center_role_assignment" {
  zone_id = "z-xxxxxx"
  principal_id = "u-xxxxxx"
  principal_type = "User"
  target_uin = "xxxxxx"
  target_type = "MemberUin"
  role_configuration_id = "rc-xxxxxx"
}
```

Import

organization identity_center_role_assignment can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_role_assignment.identity_center_role_assignment {zoneId}#{roleConfigurationId}#{targetType}#{targetUinString}#{principalType}#{principalId}
```
