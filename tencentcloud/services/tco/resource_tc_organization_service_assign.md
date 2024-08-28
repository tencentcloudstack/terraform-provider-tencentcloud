Provides a resource to create a organization service assign

Example Usage

If `management_scope` is `1`

```hcl
resource "tencentcloud_organization_service_assign" "example" {
  service_id       = 15
  management_scope = 1
  member_uins = [100037235241, 100033738111]
}
```

If `management_scope` is `2`

```hcl
resource "tencentcloud_organization_service_assign" "example" {
  service_id       = 15
  management_scope = 2
  member_uins = [100013415241, 100078908111]
  management_scope_uins = [100019287759, 100020537485]
  management_scope_node_ids = [2024256, 2024259]
}
```

Import

organization service assign can be imported using the id, e.g.
```
$ terraform import tencentcloud_organization_service_assign.example 15
```
