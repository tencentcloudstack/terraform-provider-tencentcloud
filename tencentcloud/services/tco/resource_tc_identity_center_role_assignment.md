Provides a resource to create a Organization identity center role assignment

Example Usage

```hcl
resource "tencentcloud_identity_center_role_assignment" "example" {
  zone_id               = "z-1os7c9znogct"
  principal_id          = "u-lyfm8b7qoi5l"
  principal_type        = "User"
  target_uin            = "100043911945"
  target_type           = "MemberUin"
  role_configuration_id = "rc-ihogrs0e6ceg"
}
```

Import

Organization identity center role assignment can be imported using the {zoneId}#{roleConfigurationId}#{targetType}#{targetUinString}#{principalType}, e.g.

```
terraform import tencentcloud_identity_center_role_assignment.example z-1os7c9znogct#rc-ihogrs0e6ceg#MemberUin#100043911945#User#u-lyfm8b7qoi5l
```
