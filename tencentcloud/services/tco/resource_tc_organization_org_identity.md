Provides a resource to create a organization org_identity

Example Usage

```hcl
resource "tencentcloud_organization_org_identity" "org_identity" {
  identity_alias_name = "example-iac-test"
  identity_policy {
    policy_id = 1
    policy_name = "AdministratorAccess"
    policy_type = 2
  }
  description = "iac-test"
}
```

Import

organization org_identity can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_identity.org_identity org_identity_id
```