Provides a resource to create a organization org_manage_policy

Example Usage

```hcl
resource "tencentcloud_organization_org_manage_policy" "org_manage_policy" {
  name = "FullAccessPolicy"
  content = "{\"version\":\"2.0\",\"statement\":[{\"effect\":\"allow\",\"action\":\"*\",\"resource\":\"*\"}]}"
  type = "SERVICE_CONTROL_POLICY"
  description = "Full access policy"
}
```

Import

organization org_manage_policy can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_manage_policy.org_manage_policy policy_id#type
```
