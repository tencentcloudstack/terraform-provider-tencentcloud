Provides a resource to create a organization identity center role configuration permission policy attachment

Example Usage

```hcl
resource "tencentcloud_identity_center_role_configuration_permission_policy_attachment" "example" {
  zone_id               = "z-xxxxxx"
  role_configuration_id = "rc-xxxxxx"
  role_policy_id        = 12345
}
```

Import

organization identity center role configuration permission policy attachment can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_role_configuration_permission_policy_attachment.example ${zoneId}#${roleConfigurationId}#${rolePolicyIdString}
```
