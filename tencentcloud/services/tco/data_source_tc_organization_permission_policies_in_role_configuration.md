Use this data source to query detailed information of Organization permission policies in role configuration

Example Usage

Query all permission policies in a role configuration

```hcl
data "tencentcloud_organization_permission_policies_in_role_configuration" "example" {
  zone_id               = "z-1os7c9znogct"
  role_configuration_id = "rc-ihogrs0e6ceg"
}
```

Query permission policies filtered by policy type

```hcl
data "tencentcloud_organization_permission_policies_in_role_configuration" "example" {
  zone_id               = "z-1os7c9znogct"
  role_configuration_id = "rc-ihogrs0e6ceg"
  role_policy_type      = "System"
}
```

Query permission policies filtered by policy name

```hcl
data "tencentcloud_organization_permission_policies_in_role_configuration" "example" {
  zone_id               = "z-1os7c9znogct"
  role_configuration_id = "rc-ihogrs0e6ceg"
  filter                = "AdministratorAccess"
}
```
