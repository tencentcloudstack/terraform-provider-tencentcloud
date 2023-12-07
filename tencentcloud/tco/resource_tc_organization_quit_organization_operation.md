Provides a resource to create a organization quit_organization_operation

Example Usage

```hcl
resource "tencentcloud_organization_quit_organization_operation" "quit_organization_operation" {
  org_id = 45155
}
```

Import

organization quit_organization_operation can be imported using the id, e.g.

```
terraform import tencentcloud_organization_quit_organization_operation.quit_organization_operation quit_organization_operation_id
```