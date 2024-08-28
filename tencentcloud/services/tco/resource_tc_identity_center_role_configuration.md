Provides a resource to create a organization identity_center_role_configuration

Example Usage

```hcl
resource "tencentcloud_identity_center_role_configuration" "identity_center_role_configuration" {
    zone_id = "z-xxxxxx"
    role_configuration_name = "tf-test"
    description = "test"
}
```

Import

organization identity_center_role_configuration can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_role_configuration.identity_center_role_configuration ${zoneId}#${roleConfigurationId}
```
