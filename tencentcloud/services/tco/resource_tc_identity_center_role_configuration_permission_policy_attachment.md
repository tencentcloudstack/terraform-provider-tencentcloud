Provides a resource to create a organization identity_center_role_configuration_permission_policy_attachment

Example Usage

```hcl
resource "tencentcloud_identity_center_role_configuration_permission_policy_attachment" "identity_center_role_configuration_permission_policy_attachment" {
    zone_id = "z-xxxxxx"
    role_configuration_id = "rc-xxxxxx"
    role_policy_id = xxxxxx
}
```

Import

organization identity_center_role_configuration_permission_policy_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_role_configuration_permission_policy_attachment.identity_center_role_configuration_permission_policy_attachment ${zoneId}#${roleConfigurationId}#${rolePolicyIdString}
```
