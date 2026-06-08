Use this data source to query detailed information of Organization permission policies in role configuration

Example Usage

Query all permission policies in a role configuration

```hcl
data "tencentcloud_organization_permission_policies_in_role_configuration" "example" {
  zone_id               = "z-xxxxxx"
  role_configuration_id = "rc-xxxxxx"
}
```

Query permission policies filtered by policy type

```hcl
data "tencentcloud_organization_permission_policies_in_role_configuration" "example" {
  zone_id               = "z-xxxxxx"
  role_configuration_id = "rc-xxxxxx"
  role_policy_type      = "System"
}
```

Query permission policies filtered by policy name

```hcl
data "tencentcloud_organization_permission_policies_in_role_configuration" "example" {
  zone_id               = "z-xxxxxx"
  role_configuration_id = "rc-xxxxxx"
  filter                = "AdministratorAccess"
}
```
