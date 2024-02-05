Provides a resource to create a organization org_share_unit

Example Usage

```hcl
resource "tencentcloud_organization_org_share_unit" "org_share_unit" {
  name = "iac-test"
  area = "ap-guangzhou"
  description = "iac-test"
}
```

Import

organization org_share_unit can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_share_unit.org_share_unit area#unit_id
```
