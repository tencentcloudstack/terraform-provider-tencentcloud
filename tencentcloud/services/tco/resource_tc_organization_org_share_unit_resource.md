Provides a resource to create a organization organization_org_share_unit_resource

Example Usage

```hcl
resource "tencentcloud_organization_org_share_unit_resource" "organization_org_share_unit_resource" {
  unit_id = "xxxxxx"
  area = "ap-guangzhou"
  type = "secret"
  product_resource_id = "xxxxxx"
}
```

Import

organization organization_org_share_unit_resource can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_share_unit_resource.organization_org_share_unit_resource ${unit_id}#${area}#${share_resource_type}#${product_resource_id}
```
