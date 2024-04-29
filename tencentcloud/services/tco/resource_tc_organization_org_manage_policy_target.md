Provides a resource to create a organization org_manage_policy_target

Example Usage

```hcl
resource "tencentcloud_organization_org_manage_policy_target" "org_manage_policy_target" {
  target_id = 10001
  target_type = "NODE"
  policy_id = 100001
  policy_type = "SERVICE_CONTROL_POLICY"
}
```

Import

organization org_manage_policy_target can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_manage_policy_target.org_manage_policy_target policy_type#policy_id#target_type#target_id
```
