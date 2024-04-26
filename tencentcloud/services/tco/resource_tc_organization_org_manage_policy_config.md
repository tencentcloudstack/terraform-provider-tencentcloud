Provides a resource to create a organization org_manage_policy_config

Example Usage

```hcl
resource "tencentcloud_organization_org_manage_policy_config" "org_manage_policy_config" {
  organization_id = 80001
  policy_type = "SERVICE_CONTROL_POLICY"
}
```

Import

organization org_manage_policy_config can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_manage_policy_config.org_manage_policy_config organization_id#policy_type
```
